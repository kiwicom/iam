package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtectedRoutes(t *testing.T) {
	srv := newServer()

	tests := map[string]int{
		"/":            http.StatusOK,
		"/healthcheck": http.StatusOK,
		"/teams":       http.StatusUnauthorized,
	}

	for route, code := range tests {
		req, err := http.NewRequest("GET", route, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		assert.Equal(t, code, w.Code)
	}
}
