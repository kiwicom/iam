package rest

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"

	"github.com/kiwicom/iam/internal/services/okta"
	"github.com/kiwicom/iam/internal/storage"
)

type groupsGetter interface {
	GetGroups() ([]okta.Group, error)
}

func (s *Server) handleGroupsGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := s.OktaService.GetGroups()
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
