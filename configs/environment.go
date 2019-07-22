package cfg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// InitEnv reads and set default values for environment variables.
func InitEnv() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")

	configErr := viper.ReadInConfig()
	if configErr != nil {
		log.Println("Config file failed to load. Defaulting to env.")
	}

	// viper.Unmarshal doesn't retrieve environment variables, unless they have a
	// default value, or they are specified on .env.yaml. So to make sure all envs
	// are retrieved, we set all defaults here.
	setDefaults()
}

func setDefaults() {
	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
}

// LoadConfigs loads environment variables into provided configStructs pointers.
//
// WARNING: configStructs should be pointers!
func LoadConfigs(configStructs ...interface{}) error {
	for _, c := range configStructs {
		err := viper.Unmarshal(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// ServiceConfig stores configuration values for the IAM service.
type ServiceConfig struct {
	Port         string `mapstructure:"PORT"`
	GRPCPort     string `mapstructure:"GRPC_PORT"`
	UseLocalhost bool   `mapstructure:"USE_LOCALHOST"`
	Environment  string `mapstructure:"APP_ENV"`
	Release      string `mapstructure:"SENTRY_RELEASE"`
}

// OktaConfig stores configuration values for Okta client
type OktaConfig struct {
	URL string `mapstructure:"OKTA_URL"`
}

// StorageConfig stores configuration values for storage client.
type StorageConfig struct {
	RedisHost      string        `mapstructure:"REDIS_HOST"`
	RedisPort      string        `mapstructure:"REDIS_PORT"`
	LockRetryDelay time.Duration `mapstructure:"REDIS_LOCK_RETRY_DELAY"`
	LockExpiration time.Duration `mapstructure:"REDIS_LOCK_EXPIRATION"`
}

// DatadogConfig stores configuration values for Datadog client
type DatadogConfig struct {
	Environment string `mapstructure:"APP_ENV"`
	URL         string `mapstructure:"DATADOG_ADDR"`
	AgentHost   string `mapstructure:"DD_AGENT_HOST"`
}

// SentryConfig stores configuration values for Sentry client
type SentryConfig struct {
	Token       string `mapstructure:"SENTRY_DSN"`
	Environment string `mapstructure:"APP_ENV"`
	Release     string `mapstructure:"SENTRY_RELEASE"`
}

// VaultConfig stores configuration values for Vault client
type VaultConfig struct {
	Token     string `mapstructure:"VAULT_TOKEN"`
	Address   string `mapstructure:"VAULT_ADDR"`
	Namespace string `mapstructure:"VAULT_NAMESPACE"`
}

var defaultValues = map[string]interface{}{
	"PORT":       "8080",
	"GRPC_PORT":  "8090",
	"SERVE_PATH": "/",
	// Environment used for sentry, user agent, datadog. Removes user syncing if set to dev.
	"APP_ENV": "",
	// Uses localhost intead of 0.0.0.0, useful for OSX.
	"USE_LOCALHOST": false,
	// The OKTA token and URL are only used locally, when deployed,
	// IAM fetches the token from Vault.
	"OKTA_TOKEN":             "",
	"OKTA_URL":               "",
	"REDIS_HOST":             "localhost",
	"REDIS_PORT":             "6379",
	"REDIS_LOCK_RETRY_DELAY": "1s",
	"REDIS_LOCK_EXPIRATION":  "5s",
	"SENTRY_DSN":             "",
	// The SENTRY_RELEASE value should NEVER be set manually. It's generated during docker build,
	// and it's used to track the version of the app. Useful for user agent
	// generation, and for finding regressions on Sentry.
	"SENTRY_RELEASE": "",
	// Env is taken from APP_ENV.
	"DATADOG_ADDR":    "",
	"DD_AGENT_HOST":   "",
	"VAULT_ADDR":      "",
	"VAULT_TOKEN":     "",
	"VAULT_NAMESPACE": "",
}
