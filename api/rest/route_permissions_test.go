package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/services/okta"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type permissionService struct {
	mock.Mock
}

func (u *permissionService) GetServicesPermissions(services []string) (map[string]okta.Permissions, error) {
	argsToReturn := u.Called(services)
	return argsToReturn.Get(0).(map[string]okta.Permissions), argsToReturn.Error(1)
}

func mockPermissionsRoute() (*httprouter.Router, *permissionService) {
	s := &permissionService{}
	tracer, _ := monitoring.CreateNewTracingService(monitoring.TracerOptions{
		ServiceName: "kiwi-iam",
		Environment: "test",
		Port:        "8126",
		Host:        "test",
	})

	router := httprouter.New()
	router.GET("/", getServicesPermissions(s, tracer))
	return router, s
}

func mockPermissionsRequest(router http.Handler, method, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, path, nil)
	router.ServeHTTP(response, request)
	return response
}

func TestPermissions(t *testing.T) {
	router, s := mockPermissionsRoute()

	// Unhappy paths
	response := mockPermissionsRequest(router, "GET", "/")
	assert.Equal(t, 400, response.Code, "Returns 400 when entering no service")
	assert.Equal(t, "missing services\n", response.Body.String())
	s.AssertNumberOfCalls(t, "GetServicesPermissions", 0)

	response = mockPermissionsRequest(router, "GET", "/?services=test&email=invalidemail")
	assert.Equal(t, 400, response.Code, "Returns 400 when entering an invalid email")
	assert.Equal(t, "invalid email\n", response.Body.String())
	s.AssertNumberOfCalls(t, "GetServicesPermissions", 0)

	testPermissions := map[string]okta.Permissions{
		"test": {
			"access": []string{"user1", "user2"},
			"admin":  []string{"user1"},
		},
	}

	// Happy path
	s.On("GetServicesPermissions", []string{"test"}).Return(testPermissions, nil)

	response = mockPermissionsRequest(router, "GET", "/?services=test")
	assert.Equal(t, 200, response.Code, "Returns 200 on success")

	responsePermissions := make(map[string]okta.Permissions)
	_ = json.Unmarshal(response.Body.Bytes(), &responsePermissions)

	assert.Equal(t, testPermissions, responsePermissions)
	s.AssertNumberOfCalls(t, "GetServicesPermissions", 1)
}
