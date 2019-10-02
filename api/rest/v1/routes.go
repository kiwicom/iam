package rest

import "github.com/gorilla/mux"

// routes handles registering all routes. All routes should be added here.
func (s *server) routes() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/", s.handleHello())
}
