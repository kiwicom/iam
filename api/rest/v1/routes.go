package rest

import "github.com/gorilla/mux"

// routes handles registering all routes. All routes should be added here.
func (s *Server) routes() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.handleHello())
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck())
	s.router.HandleFunc("/teams", s.middlewareSecurity(s.handleTeamsGET()))
}
