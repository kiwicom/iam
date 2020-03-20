package rest

import (
	"net/http"
	"strings"

	tracingRouter "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// routes handles registering all routes. All routes should be added here.
func (s *Server) routes() {
	const wellKnownFolder string = ".well-known"

	s.Router = tracingRouter.NewRouter(tracingRouter.WithServiceName(s.ServiceName))

	s.Router.HandleFunc("/", s.handleHello())
	s.Router.HandleFunc("/healthcheck", s.handleHealthcheck())
	s.Router.HandleFunc("/v1/user", s.middlewareSecurity(s.handleUserGET()))
	s.Router.HandleFunc("/v1/groups", s.middlewareSecurity(s.handleGroupsGET()))

	s.Router.PathPrefix("/" + wellKnownFolder + "/").Handler(DisableDirectoryListingHandler(
		http.StripPrefix("/"+wellKnownFolder+"/", http.FileServer(http.Dir(wellKnownFolder))),
	))
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
