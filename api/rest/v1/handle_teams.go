package rest

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
)

func (s *Server) handleTeamsGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := s.OktaService.GetTeams()
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(teams); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
