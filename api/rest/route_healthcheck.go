package rest

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// healthcheck is a route for healthcheck
func healthcheck(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	var message = "Ok"
	_, err := w.Write([]byte(message))

	if err != nil {
		log.Println("[ERROR] Failed to return healthcheck: ", err)
	}
}
