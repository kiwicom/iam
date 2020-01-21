package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"

	grpcAPI "github.com/kiwicom/iam/api/grpc"
	pb "github.com/kiwicom/iam/api/grpc/v1"
	restAPI "github.com/kiwicom/iam/api/rest/v1"
	cfg "github.com/kiwicom/iam/configs"
	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/security/secrets"
	"github.com/kiwicom/iam/internal/services/okta"
	"github.com/kiwicom/iam/internal/storage"

	// AppEngine is added manually due to issues with go.mod
	// https://skypicker.slack.com/archives/CA154LA5T/p1560781760024700
	_ "google.golang.org/appengine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func syncOkta(client *okta.Client) {
	log.Println("Start caching users")
	client.SyncUsers()
	log.Println("Start caching groups")
	client.SyncGroups()

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		client.SyncUsers()
		log.Println("Start caching groups", tick.Round(time.Second))
		client.SyncGroups()
	}
}

func syncVault(client *secrets.VaultManager) {
	log.Println("Start token sync with Vault")
	err := client.SyncAppTokens()

	if err != nil {
		log.Println("[ERROR] failed to sync tokens with Vault: ", err)
		raven.CaptureError(err, nil)
	}

	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching app tokens from Vault", tick.Round(time.Second))
		err := client.SyncAppTokens()

		if err != nil {
			log.Println("[ERROR] failed to sync tokens with Vault: ", err)
			raven.CaptureError(err, nil)
		}
	}
}

// Cacher contains methods needed from a cache
type Cacher interface {
	Del(key string) error
}

func clearGroupsLastSync(cache Cacher) {
	log.Println("Clearing GroupLastSync")
	err := cache.Del("groups-sync-timestamp")
	if err != nil {
		log.Println("[ERROR] failed to delete GroupsLastSync from cache: ", err)
		raven.CaptureError(err, nil)
	}

	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Clearing GroupLastSync", tick.Round(time.Second))
		err = cache.Del("groups-sync-timestamp")
		if err != nil {
			log.Println("[ERROR] failed to delete GroupsLastSync from cache: ", err)
			raven.CaptureError(err, nil)
		}
	}
}

func createSecretManager(vault cfg.VaultConfig) secrets.SecretManager {
	// Load data from Vault and set them if possible
	vaultClient, vaultErr := secrets.CreateNewVaultClient(
		vault.Address,
		vault.Token,
		vault.Namespace,
	)
	if vaultErr == nil {
		// This sync needs to happen synchronously
		err := vaultClient.SyncAppSettings()

		// If Vault is set up, but connection fails or settings are empty then kill the app as oktaToken will not be available
		if err != nil {
			panic(err)
		}

		go capturePanic(func() { syncVault(vaultClient) })

		return vaultClient
	}

	log.Println("Vault integration disabled: ", vaultErr)

	localSecretManager := secrets.CreateNewLocalSecretManager()

	return localSecretManager
}

func initErrorTracking(sentry cfg.SentryConfig) {
	if sentry.Token == "" {
		log.Println("SENTRY_DSN is not set. Error logging disabled.")
		return
	}
	err := raven.SetDSN(sentry.Token)
	if err != nil {
		log.Println("[ERROR] Failed to set Raven DSN: ", err)
	}

	raven.SetEnvironment(sentry.Environment)
	raven.SetRelease(sentry.Release)
}

func capturePanic(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)

			if e, ok := err.(error); ok {
				raven.CaptureError(e, nil)
			}

			if m, ok := err.(string); ok {
				raven.CaptureMessage(m, nil)
			}
		}
	}()

	f()
}

//nolint:funlen // might need some refactoring in the future but is otherwise is readable
func main() {
	cfg.InitEnv()

	var (
		iamConfig     cfg.ServiceConfig
		oktaConfig    cfg.OktaConfig
		storageConfig cfg.StorageConfig
		datadogConfig cfg.DatadogConfig
		sentryConfig  cfg.SentryConfig
		vaultConfig   cfg.VaultConfig
	)

	// If there is an error loading the envs kill the app, as nothing will work without them.
	if err := cfg.LoadConfigs(&iamConfig, &oktaConfig, &storageConfig, &datadogConfig, &sentryConfig, &vaultConfig); err != nil {
		log.Println("[ERROR]", err.Error())
		panic(err)
	}

	initErrorTracking(sentryConfig)
	secretManager := createSecretManager(vaultConfig)

	// Datadog tracer
	tracer, _ := monitoring.CreateNewTracingService(monitoring.TracerOptions{
		ServiceName: "kiwi-iam",
		Environment: iamConfig.Environment,
		Port:        "8126",
		Host:        datadogConfig.AgentHost,
	})

	// Metrics initialization
	metricClient, metricErr := monitoring.CreateNewMetricService(monitoring.MetricSettings{
		Host:        datadogConfig.AgentHost,
		Port:        "8125",
		Namespace:   "kiwi_iam.",
		Environment: iamConfig.Environment,
	})
	if metricErr != nil {
		log.Println("[ERROR]", metricErr)
	}

	cache := storage.NewRedisCache(
		storageConfig.RedisHost,
		storageConfig.RedisPort,
		1,
	)
	lock := storage.NewLockManager(
		cache,
		storageConfig.LockRetryDelay,
		storageConfig.LockExpiration,
	)
	oktaToken, _ := secretManager.GetSetting("OKTA_TOKEN")
	oktaClient := okta.NewClient(&okta.ClientOpts{
		BaseURL:     oktaConfig.URL,
		AuthToken:   oktaToken,
		Cache:       cache,
		LockManager: lock,
		IAMConfig:   &iamConfig,
		Metrics:     metricClient,
	})

	restServer := restAPI.NewServer("kiwi-iam.http.router")
	restServer.OktaService = oktaClient
	restServer.SecretManager = secretManager
	restServer.MetricClient = metricClient
	restServer.Tracer = tracer

	// 0.0.0.0 is specified to allow listening in Docker
	var address = "0.0.0.0"
	if iamConfig.UseLocalhost {
		address = "localhost"
	}

	var serveAddr = net.JoinHostPort(address, iamConfig.Port)
	server := &http.Server{
		Handler:      restServer.Router,
		Addr:         serveAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if iamConfig.Environment != "dev" {
		go capturePanic(func() { clearGroupsLastSync(cache) })
		go capturePanic(func() { syncOkta(oktaClient) })
	}

	log.Println("🚀 REST server starting on " + serveAddr)
	go capturePanic(func() { _ = server.ListenAndServe() })

	// GRPC Init

	// Create listener
	grpcAddress := net.JoinHostPort(address, iamConfig.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s, _ := grpcAPI.CreateServer(oktaClient)

	creds, err := credentials.NewServerTLSFromFile(iamConfig.GRPCCertFile, iamConfig.GRPCKeyFile)
	if err != nil {
		log.Println("TLS disabled on GRPC:", err)
		raven.CaptureError(err, nil)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcAPI.UnarySecurityWrapper(secretManager)), grpc.Creds(creds))
	reflection.Register(grpcServer)

	pb.RegisterKiwiIAMAPIServer(grpcServer, s)

	log.Printf("🚀 GRPC server listening on %s", grpcAddress)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
