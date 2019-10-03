package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	s := Server{}

	req, err := http.NewRequest("GET", "/healthcheck", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	handler := s.handleHealthcheck()

	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body := w.Body.String()
	assert.Equal(t, "Ok", body)
}
