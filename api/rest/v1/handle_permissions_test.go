package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiwicom/iam/api"
	"github.com/kiwicom/iam/internal/services/okta"
)

// nolint:unparam // even though method is currently always "GET" we might decide to use other methods in the future
func mockPermissionsRequest(handler http.HandlerFunc, method, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, path, nil)
	handler.ServeHTTP(response, request)
	return response
}

func TestServicePermissions(t *testing.T) {
	userService := &api.MockOktaService{}
	server := setupServer()
	server.OktaService = userService

	// Unhappy paths
	response := mockPermissionsRequest(server.handlePermissionsGET(), "GET", "/")
	assert.Equal(t, 400, response.Code, "Returns 400 when entering no service")
	assert.Equal(t, "missing services\n", response.Body.String())
	userService.AssertNumberOfCalls(t, "GetServicesPermissions", 0)

	response = mockPermissionsRequest(server.handlePermissionsGET(), "GET", "/?services=test&email=invalidemail")
	assert.Equal(t, 400, response.Code, "Returns 400 when entering an invalid email")
	assert.Equal(t, "invalid email\n", response.Body.String())
	userService.AssertNumberOfCalls(t, "GetServicesPermissions", 0)

	expectedPermissions := map[string]okta.Permissions{
		"test": {
			"access": []string{"user1", "user2"},
			"admin":  []string{"user1"},
		},
	}

	// Happy path
	userService.On("GetServicesPermissions", []string{"test"}).Return(expectedPermissions, nil)

	response = mockPermissionsRequest(server.handlePermissionsGET(), "GET", "/?services=test")
	assert.Equal(t, 200, response.Code, "Returns 200 on success")

	actualPermissions := make(map[string]okta.Permissions)
	_ = json.Unmarshal(response.Body.Bytes(), &actualPermissions)

	assert.Equal(t, expectedPermissions, actualPermissions)
	userService.AssertNumberOfCalls(t, "GetServicesPermissions", 1)
	userService.AssertNumberOfCalls(t, "GetUserPermissions", 0)
}

func TestUserPermissions(t *testing.T) {
	userService := &api.MockOktaService{}
	server := setupServer()
	server.OktaService = userService

	expectedPermissions := map[string][]string{"test": {"access", "admin"}}
	userService.On("GetUserPermissions", "user@test.com", []string{"test"}).Return(expectedPermissions, nil)

	response := mockPermissionsRequest(server.handlePermissionsGET(), "GET", "/?services=test&email=user@test.com")
	assert.Equal(t, 200, response.Code, "Returns 200 on success")

	actualPermissions := make(map[string][]string)
	_ = json.Unmarshal(response.Body.Bytes(), &actualPermissions)

	assert.Equal(t, expectedPermissions, actualPermissions)
	userService.AssertNumberOfCalls(t, "GetUserPermissions", 1)
	userService.AssertNumberOfCalls(t, "GetServicePermissions", 0)
}
