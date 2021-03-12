package rest

import (
	"log"
	"net/http"
)

// handleHealthcheck is a route for healthcheck.
func (s *Server) handleHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "Ok"
		_, err := w.Write([]byte(message))
		if err != nil {
			log.Println("[ERROR] Failed to return healthcheck: ", err)
		}
	}
}
