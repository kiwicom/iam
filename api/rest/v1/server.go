package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// server houses all dependencies and routing of the server
type server struct {
	router *mux.Router
}

func newServer() server {
	s := server{}

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
