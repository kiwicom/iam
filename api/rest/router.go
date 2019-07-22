package rest

import (
	"net/http"
	"strings"

	"github.com/iam/internal/monitoring"
	"github.com/iam/internal/security/secrets"
	"github.com/iam/internal/services/okta"
	"github.com/julienschmidt/httprouter"

	tracingRouter "gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

const wellKnownFolder string = ".well-known"

type metricService interface {
	// Incr increments by 1 a metric identified by name.
	// tags should be in format name:value and can be created with Tag function to escape the values
	Incr(string, ...string)
}

// CreateRouter creates a new router instance
func CreateRouter(
	serviceName string,
	oktaClient *okta.Client,
	secretManager secrets.SecretManager,
	metricClient metricService,
	tracer *monitoring.Tracer) *tracingRouter.Router {
	router := tracingRouter.New(tracingRouter.WithServiceName(serviceName))

	router.Handler(
		"GET",
		"/"+wellKnownFolder+"/*filepath",
		DisableDirectoryListingHandler(
			http.StripPrefix("/"+wellKnownFolder+"/", http.FileServer(http.Dir(wellKnownFolder))),
		),
	)

	// Healthcheck routes.
	router.GET("/healthcheck", healthcheck)

	// Hello World Route
	router.GET("/", sayHello)

	addEndpoint := func(path string, handler httprouter.Handle) {
		router.GET(path,
			AuthWrapper(
				handler,
				secretManager,
				metricClient,
			),
		)
	}

	// App routes
	addEndpoint("/v1/user", getOktaUserByEmail(oktaClient, tracer))

	addEndpoint("/v1/teams", getTeams(oktaClient))

	addEndpoint("/v1/groups", getGroups(oktaClient))

	router.PanicHandler = panicHandler

	return router
}

// DisableDirectoryListingHandler prevents directory listings to be returned to the user-agent.
func DisableDirectoryListingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
