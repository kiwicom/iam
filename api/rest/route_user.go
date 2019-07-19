package rest

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/getsentry/raven-go"
	"github.com/iam/monitoring"
	"github.com/iam/security"
	"github.com/iam/services/okta"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type userDataService interface {
	GetUser(string) (okta.User, error)
	AddPermissions(*okta.User, string) error
}

// getOktaUserByEmail looks up an Okta user by email
func getOktaUserByEmail(client userDataService, tracer *monitoring.Tracer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		params, paramErr := validateUsersParams(r.URL.RawQuery)
		if paramErr != nil {
			http.Error(w, paramErr.Error(), http.StatusBadRequest)
			return
		}
		email := params["email"]
		permissions := params["permissions"] == "true"

		service, getServiceErr := security.GetService(r.Header.Get("User-Agent"))
		if getServiceErr != nil {
			http.Error(w, "Invalid user agent", http.StatusBadRequest)
			return
		}

		// getUser just wraps GetUser in tracing
		getUser := func() (okta.User, error) {
			span, _ := tracer.StartSpanWithContext(r.Context(), "user-data", "okta-controller", "http")
			defer tracer.FinishSpan(span)
			oktaUser, err := client.GetUser(email)

			return oktaUser, err
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
			span, _ := tracer.StartSpanWithContext(r.Context(), "permissions", "okta-controller", "http")
			defer tracer.FinishSpan(span)
			permErr := client.AddPermissions(&oktaUser, service.Name)

			return permErr
		}

		if permissions {
			permErr := addPermissions()
			if permErr != nil {
				log.Println("[ERROR]", permErr.Error())
				raven.CaptureError(permErr, nil)
			}
		}
		oktaUser.OktaID = ""           // OktaID is used only internally
		oktaUser.GroupMembership = nil // GroupMembership is used only internally

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)

		mapUser, err := formatUser(&oktaUser, permissions)
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
		"email":       values.Get("email"),
		"permissions": values.Get("permissions"),
	}

	if params["email"] == "" {
		return nil, errors.New("missing email")
	}
	if _, err := mail.ParseAddress(params["email"]); err != nil {
		return nil, errors.New("invalid email")
	}

	return params, nil
}

// formatUser converts the given user to map and keeps or removes the
// permissions field based on the withPermissions parameter.
func formatUser(s *okta.User, withPermissions bool) (map[string]interface{}, error) {
	str, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(str, &m)
	if !withPermissions {
		delete(m, "permissions")
	}
	return m, err
}
