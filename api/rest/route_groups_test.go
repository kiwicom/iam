package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iam/services/okta"
	"github.com/iam/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGroupsGetter struct {
	mock.Mock
}

func (g *mockGroupsGetter) GetGroups() ([]okta.Group, error) {
	args := g.Called()
	return args.Get(0).([]okta.Group), args.Error(1)
}

func TestGetGroupsErrors(t *testing.T) {
	g := &mockGroupsGetter{}

	request, _ := http.NewRequest("GET", "/", nil)
	router := httprouter.New()
	router.GET("/", getGroups(g))

	// Generic error
	errMessage := "internal error that shouldn't be exposed"
	g.On("GetGroups").Return([]okta.Group{}, errors.New(errMessage)).Once()
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code)
	assert.NotEqual(t, errMessage, response.Body.String())

	// No value found in cache
	g.On("GetGroups").Return([]okta.Group{}, storage.ErrNotFound).Once()
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, 503, response.Code)
	assert.NotEqual(t, "", response.Header().Get("Retry-After"))
	g.AssertExpectations(t)
}

func TestGetGroups(t *testing.T) {
	g := &mockGroupsGetter{}

	request, _ := http.NewRequest("GET", "/", nil)
	router := httprouter.New()
	router.GET("/", getGroups(g))

	groups := []okta.Group{{ID: "id1", Name: "Group 1"}}
	g.On("GetGroups").Return(groups, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	expected := "[{\"id\":\"id1\",\"name\":\"Group 1\",\"description\":\"\",\"lastMembershipUpdated\":\"0001-01-01T00:00:00Z\"}]\n"
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, response.Body.String())
	g.AssertExpectations(t)
}
