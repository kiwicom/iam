package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiwicom/iam/internal/monitoring"
	"github.com/kiwicom/iam/internal/services/okta"
)

var testUser = okta.User{
	FirstName:   "Test",
	LastName:    "Tester",
	Position:    "QA Tester",
	Permissions: []string{"action:read"},
}

func setupServer() *Server {
	s := &Server{}
	tracer, _ := monitoring.CreateNewTracingService(monitoring.TracerOptions{
		ServiceName: "kiwi-iam",
		Environment: "test",
		Port:        "8126",
		Host:        "test",
	})

	s.tracer = tracer

	return s
}

func TestMissingQuery(t *testing.T) {
	userService := &mockOktaService{}
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()

	handler.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when entering wrong email")

	responseBody := response.Body.String()
	assert.Equal(t, "missing email\n", responseBody, "Returns correct body")
	userService.AssertNotCalled(t, "GetUser")
	userService.AssertNotCalled(t, "AddPermissions")
}

func TestWrongEmail(t *testing.T) {
	userService := &mockOktaService{}
	request, _ := http.NewRequest("GET", "/?email=testest", nil)
	response := httptest.NewRecorder()
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()

	handler.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when entering wrong email")

	responseBody := response.Body.String()
	assert.Equal(t, "invalid email\n", responseBody, "Returns correct body")
	userService.AssertNotCalled(t, "GetUser")
	userService.AssertNotCalled(t, "AddPermissions")
}

func TestMissingUserAgent(t *testing.T) {
	userService := &mockOktaService{}
	request, _ := http.NewRequest("GET", "/?email=test@test.com", nil)
	response := httptest.NewRecorder()
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()
	handler.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when a user agent header is missing")

	responseBody := response.Body.String()
	assert.Equal(t, "Invalid user agent\n", responseBody, "Returns correct body")
	userService.AssertNotCalled(t, "GetUser")
	userService.AssertNotCalled(t, "AddPermissions")
}

func TestHappyPathWithPermissions(t *testing.T) {
	// Success response
	userService := &mockOktaService{}
	request, _ := http.NewRequest("GET", "/?email=test@test.com&permissions=true", nil)
	request.Header.Set("User-Agent", "service/0 (Kiwi.com test)")
	response := httptest.NewRecorder()
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()
	userService.On("GetUser", "test@test.com").Return(testUser, nil)
	userService.On("AddPermissions", &testUser, "service").Return(nil)

	handler.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Returns 200 on success")

	responseJSON := response.Body.Bytes()
	var responseUser okta.User
	_ = json.Unmarshal(responseJSON, &responseUser)

	// For some reason response adds a extra line break
	assert.Equal(t, testUser, responseUser, "Returns correct body")
	userService.AssertNumberOfCalls(t, "GetUser", 1)
	userService.AssertNumberOfCalls(t, "AddPermissions", 1)
}

func TestHappyPathNoPermissions(t *testing.T) {
	// Success response
	urls := []string{
		"/?email=test@test.com&permissions=false",
		"/?email=test@test.com", // default value of permissions is false
	}

	for _, url := range urls {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("User-Agent", "service/0 (Kiwi.com test)")
		response := httptest.NewRecorder()

		userService := &mockOktaService{}
		server := setupServer()
		server.oktaService = userService

		handler := server.handleUserGET()
		userService.On("GetUser", "test@test.com").Return(testUser, nil)
		userService.On("AddPermissions", &testUser, "service").Return(nil)

		handler.ServeHTTP(response, request)
		assert.Equal(t, 200, response.Code, "Returns 200 on success")

		responseJSON := response.Body.Bytes()
		var responseMap map[string]interface{}
		_ = json.Unmarshal(responseJSON, &responseMap)

		var expectedUser map[string]interface{}
		str, _ := json.Marshal(testUser)
		_ = json.Unmarshal(str, &expectedUser)
		delete(expectedUser, "permissions")

		// For some reason response adds a extra line break
		assert.Equal(t, expectedUser, responseMap, "Returns correct body")
		userService.AssertNumberOfCalls(t, "GetUser", 1)
		userService.AssertNumberOfCalls(t, "AddPermissions", 0)
	}
}

func TestControllerFailurePath(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=bs@test.com", nil)
	request.Header.Set("User-Agent", "service/0 (Kiwi.com test)")
	response := httptest.NewRecorder()

	userService := &mockOktaService{}
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()
	userService.On("GetUser", "bs@test.com").Return(okta.User{}, errors.New("boom"))
	handler.ServeHTTP(response, request)
	assert.Equal(t, 500, response.Code, "Returns 500 on controller failure")

	responseBody := response.Body.String()
	assert.Equal(t, "Service unavailable\n", responseBody, "Returns error correct body")
	userService.AssertNumberOfCalls(t, "GetUser", 1)
	userService.AssertNotCalled(t, "AddPermissions")
}

func TestNotFoundPath(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=notfound@test.com", nil)
	request.Header.Set("User-Agent", "service/0 (Kiwi.com test)")
	response := httptest.NewRecorder()

	userService := &mockOktaService{}
	server := setupServer()
	server.oktaService = userService

	handler := server.handleUserGET()
	userService.On("GetUser", "notfound@test.com").Return(okta.User{}, okta.ErrUserNotFound)
	handler.ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "Returns 404 on user not found")

	responseBody := response.Body.String()
	assert.Equal(t, "User notfound@test.com not found\n", responseBody, "Returns correct body")
	userService.AssertNumberOfCalls(t, "GetUser", 1)
	userService.AssertNotCalled(t, "AddPermissions")
}
