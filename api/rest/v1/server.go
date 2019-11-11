package rest

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	tracingRouter "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"

	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/security/secrets"
	"github.com/kiwicom/iam/internal/services/okta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type oktaService interface {
	AddPermissions(*okta.User, string) error
	GetTeams() (map[string]int, error)
	GetUser(string) (okta.User, error)
	GetServicesPermissions([]string) (map[string]okta.Permissions, error)
	GetUserPermissions(string, []string) (map[string][]string, error)
}

type metricService interface {
	// Incr increments by 1 a metric identified by name.
	// tags should be in format name:value and can be created with Tag function to escape the values
	Incr(string, ...string)
}

// Server houses all dependencies and routing of the server
type Server struct {
	Router        *tracingRouter.Router
	SecretManager secrets.SecretManager
	MetricClient  metricService
	OktaService   oktaService
	Tracer        *monitoring.Tracer
	// ServiceName is used for tracing purposes
	ServiceName string
}

// NewServer creates a new instance of server and sets up routes
func NewServer() *Server {
	s := Server{}
	s.routes()

	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
