package rest

import tracingRouter "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"

// routes handles registering all routes. All routes should be added here.
func (s *Server) routes() {
	s.Router = tracingRouter.NewRouter(tracingRouter.WithServiceName(s.ServiceName))
	s.Router.HandleFunc("/", s.handleHello())
	s.Router.HandleFunc("/healthcheck", s.handleHealthcheck())
	s.Router.HandleFunc("/v1/teams", s.middlewareSecurity(s.handleTeamsGET()))
	s.Router.HandleFunc("/v1/user", s.middlewareSecurity(s.handleUserGET()))
	s.Router.HandleFunc("/v1/permissions", s.middlewareSecurity(s.handlePermissionsGET()))
}
