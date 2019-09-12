package rest

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/julienschmidt/httprouter"

	"github.com/kiwicom/iam/internal/services/okta"
	"github.com/kiwicom/iam/internal/storage"
)

type groupsGetter interface {
	GetGroups() ([]okta.Group, error)
}

func getGroups(client groupsGetter) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		groups, err := client.GetGroups()
		if err == storage.ErrNotFound {
			// No value available for groups yet
			w.Header().Add("Retry-After", "30")
			http.Error(w, "Groups not loaded yet, try later", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(groups); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
