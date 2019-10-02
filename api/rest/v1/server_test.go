package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtectedRoutes(t *testing.T) {
	srv := server{}
	srv.routes()

	tests := map[string]int{
		"/":            200,
		"/healthcheck": 200,
	}

	for route, code := range tests {
		req, err := http.NewRequest("GET", route, nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		assert.Equal(t, code, w.Code)

	}
}
