package rest

import (
	"net/http"
)

func (s *Server) middlewareDeprecation(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Deprecated", "true")

		h(w, r)
	}
}
