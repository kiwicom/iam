package rest

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/getsentry/raven-go"

	"github.com/kiwicom/iam/internal/security"
	"github.com/kiwicom/iam/internal/services/okta"
)

// handleUser looks up an Okta user by email
func (s *Server) handleUserGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, paramErr := validateUsersParams(r.URL.RawQuery)
		if paramErr != nil {
			http.Error(w, paramErr.Error(), http.StatusBadRequest)
			return
		}
		email := params["email"]
		serviceName := params["service"]

		if serviceName == "" {
			service, getServiceErr := security.GetService(r.Header.Get("User-Agent"))
			if getServiceErr != nil {
				http.Error(w, "Missing service and invalid user agent", http.StatusBadRequest)
				return
			}

			serviceName = service.Name
		}

		// getUser just wraps GetUser in tracing
		getUser := func() (*okta.User, error) {
			span, _ := s.Tracer.StartSpanWithContext(r.Context(), "user-data", "okta-controller", "http")
			defer s.Tracer.FinishSpan(span)
			oktaUser, err := s.OktaService.GetUser(email)

			return &oktaUser, err
		}
		oktaUser, err := getUser()
		if err == okta.ErrUserNotFound {
			http.Error(w, "User "+email+" not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		// addPermissions just wraps AddPermissions with tracing
		addPermissions := func() error {
			span, _ := s.Tracer.StartSpanWithContext(r.Context(), "permissions", "okta-controller", "http")
			defer s.Tracer.FinishSpan(span)
			permErr := s.OktaService.AddPermissions(oktaUser, serviceName)

			return permErr
		}

		permErr := addPermissions()
		if permErr != nil {
			log.Println("[ERROR]", permErr.Error())
			raven.CaptureError(permErr, nil)
		}

		oktaUser.OktaID = ""           // OktaID is used only internally
		oktaUser.GroupMembership = nil // GroupMembership is used only internally

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)

		mapUser, err := formatUser(oktaUser)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if err := je.Encode(mapUser); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}

// validateUsersParams validates query parameters for the users endpoint.
func validateUsersParams(rawQuery string) (map[string]string, error) {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, errors.New("invalid query string")
	}

	params := map[string]string{
		"email":   values.Get("email"),
		"service": values.Get("service"),
	}

	if params["email"] == "" {
		return nil, errors.New("missing email")
	}
	if _, err := mail.ParseAddress(params["email"]); err != nil {
		return nil, errors.New("invalid email")
	}

	return params, nil
}

// formatUser converts the given user to map
func formatUser(s *okta.User) (map[string]interface{}, error) {
	str, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(str, &m)

	return m, err
}
