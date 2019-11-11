package rest

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/getsentry/raven-go"

	"github.com/kiwicom/iam/internal/services/okta"
)

type permissionsParams struct {
	services []string
	email    string
}

func (s *Server) handlePermissionsGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, paramErr := validatePermissionsParams(r.URL.RawQuery)
		if paramErr != nil {
			http.Error(w, paramErr.Error(), http.StatusBadRequest)
			return
		}

		var permissions interface{}
		var err error

		if params.email == "" {
			span, _ := s.Tracer.StartSpanWithContext(r.Context(), "services-permissions", "okta-controller", "http")
			defer s.Tracer.FinishSpan(span)
			permissions, err = s.OktaService.GetServicesPermissions(params.services)
		} else {
			span, _ := s.Tracer.StartSpanWithContext(r.Context(), "user-permissions", "okta-controller", "http")
			defer s.Tracer.FinishSpan(span)
			permissions, err = s.OktaService.GetUserPermissions(params.email, params.services)
		}

		if err == okta.ErrNotReady {
			w.Header().Add("Retry-After", "30")
			http.Error(w, "Data not ready", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)

		if err := je.Encode(permissions); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}

// validatePermissionsParams validates query parameters for the permissions endpoint.
func validatePermissionsParams(rawQuery string) (permissionsParams, error) {
	params := permissionsParams{}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return params, errors.New("invalid query string")
	}

	if values.Get("services") == "" {
		return params, errors.New("missing services")
	}

	params.services = strings.Split(values.Get("services"), ",")
	params.email = values.Get("email")

	if _, err := mail.ParseAddress(params.email); params.email != "" && err != nil {
		return params, errors.New("invalid email")
	}

	return params, nil
}
