package rest

import (
	"log"
	"net/http"
)

// handleHello is a route for HelloWorld
func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "42"
		_, err := w.Write([]byte(message))
		if err != nil {
			log.Println("Hello World failed: ", err)
		}
	}
}
