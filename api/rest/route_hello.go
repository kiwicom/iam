package rest

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// sayHello is a route for HelloWorld
func sayHello(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	message := "42"
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println("Hello World failed: ", err)
	}
}
