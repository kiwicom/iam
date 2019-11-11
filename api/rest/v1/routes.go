package rest

import "github.com/gorilla/mux"

// routes handles registering all routes. All routes should be added here.
func (s *Server) routes() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/", s.handleHello())
	s.Router.HandleFunc("/healthcheck", s.handleHealthcheck())
	s.Router.HandleFunc("/teams", s.middlewareSecurity(s.handleTeamsGET()))
	s.Router.HandleFunc("/user", s.middlewareSecurity(s.handleUserGET()))
}
