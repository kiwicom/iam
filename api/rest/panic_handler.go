package rest

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/iam/api"
)

// panicHandler prevents the server crashing due to unhandled panics
func panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	apiError, castingOk := err.(api.Error)
	if castingOk {
		raven.CaptureError(apiError, nil)
		http.Error(w, apiError.Message, apiError.Code)
		log.Println("[ERROR]", apiError)
		return
	}

	if errorType, castingOk := err.(error); castingOk {
		raven.CaptureError(errorType, nil)
	}
	log.Panic(err)
}
