package rest

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/julienschmidt/httprouter"

	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/services/okta"
)

type permissionDataService interface {
	GetServicesPermissions([]string) (map[string]okta.Permissions, error)
}

type permissionsParams struct {
	services []string
	email    string
}

// getPermissions looks up permissions
func getServicesPermissions(client permissionDataService, tracer *monitoring.Tracer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		params, paramErr := validatePermissionsParams(r.URL.RawQuery)
		if paramErr != nil {
			http.Error(w, paramErr.Error(), http.StatusBadRequest)
			return
		}

		// getPermissions just wraps GetServicesPermissions in tracing
		getPermissions := func() (map[string]okta.Permissions, error) {
			span, _ := tracer.StartSpanWithContext(r.Context(), "permissions-data", "okta-controller", "http")
			defer tracer.FinishSpan(span)
			permissions, err := client.GetServicesPermissions(params.services)

			return permissions, err
		}

		permissions, err := getPermissions()
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
