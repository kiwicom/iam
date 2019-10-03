package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/kiwicom/iam/internal/security/secrets"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type oktaService interface {
	GetTeams() (map[string]int, error)
}

type metricService interface {
	// Incr increments by 1 a metric identified by name.
	// tags should be in format name:value and can be created with Tag function to escape the values
	Incr(string, ...string)
}

// server houses all dependencies and routing of the server
type server struct {
	router        *mux.Router
	secretManager secrets.SecretManager
	metricClient  metricService
	oktaService   oktaService
}

func newServer() server {
	s := server{}
	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
