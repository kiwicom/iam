package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeprecationWrapper(t *testing.T) {
	s := Server{}

	handler := s.middlewareDeprecation(s.handleHello())

	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	assert.Equal(t, "true", response.Header().Get("Deprecated"), "Deprecated header should be in response")
}
