package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiwicom/iam/internal/security/secrets"
)

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
}

func newServer() server {
	s := server{}
	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
