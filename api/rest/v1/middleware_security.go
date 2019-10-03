package rest

import (
	"log"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/kiwicom/iam/api"
	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/security"
)

// AuthWrapper wraps a router to validate the authentication token
func (s *Server) middlewareSecurity(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.checkAuth(r)
		if err != nil {
			if apiErr, ok := err.(api.Error); ok {
				http.Error(w, apiErr.Message, apiErr.Code)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			log.Println("[ERROR]", err.Error())
			return
		}

		// Delegate request to the given handle
		h(w, r)
	}
}

// checkAuth checks if user has proper token + user agent
func (s *Server) checkAuth(r *http.Request) error {
	requestToken, err := security.GetToken(r.Header.Get("Authorization"))
	if err != nil {
		return api.Error{Message: "Use the Bearer {token} authorization scheme", Code: http.StatusUnauthorized}
	}
	userAgent := r.Header.Get("User-Agent")

	service, err := security.GetService(userAgent)
	if err != nil {
		return api.Error{Message: err.Error(), Code: http.StatusUnauthorized}
	}

	if span, ok := tracer.SpanFromContext(r.Context()); ok {
		span.SetTag("user-agent", userAgent)
		span.SetTag("service-name", service.Name)
	}

	tokenErr := security.VerifyToken(s.secretManager, service, requestToken)

	if tokenErr != nil {
		return api.Error{Message: "Unauthorized: " + tokenErr.Error(), Code: http.StatusUnauthorized}
	}

	s.metricClient.Incr(
		"incoming.requests",
		monitoring.Tag("service-name", service.Name),
		monitoring.Tag("service-environment", service.Environment),
	)

	return nil
}
