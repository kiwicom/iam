package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// revive:disable:confusing-naming createRouter is for testing, CreateRouter is used normally
func createRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", addDeprecationWarning(sayHello))
	return router
}

func TestDeprecationWrapper(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	createRouter().ServeHTTP(response, request)
	assert.Equal(t, "true", response.Header().Get("Deprecated"), "True is expected")
}
