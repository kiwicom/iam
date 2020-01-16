package api

import (
	"github.com/stretchr/testify/mock"

	"github.com/kiwicom/iam/internal/services/okta"
)

type MockOktaService struct {
	mock.Mock
}

func (o *MockOktaService) AddPermissions(user *okta.User, service string) error {
	argsToReturn := o.Called(user, service)
	return argsToReturn.Error(0)
}

func (o *MockOktaService) GetTeams() (map[string]int, error) {
	args := o.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (o *MockOktaService) GetUser(email string) (okta.User, error) {
	argsToReturn := o.Called(email)
	return argsToReturn.Get(0).(okta.User), argsToReturn.Error(1)
}

func (o *MockOktaService) GetServicesPermissions(services []string) (map[string]okta.Permissions, error) {
	argsToReturn := o.Called(services)
	return argsToReturn.Get(0).(map[string]okta.Permissions), argsToReturn.Error(1)
}

func (o *MockOktaService) GetUserPermissions(email string, services []string) (map[string][]string, error) {
	argsToReturn := o.Called(email, services)
	return argsToReturn.Get(0).(map[string][]string), argsToReturn.Error(1)
}

func (o *MockOktaService) GetGroups() ([]okta.Group, error) {
	argsToReturn := o.Called()
	return argsToReturn.Get(0).([]okta.Group), argsToReturn.Error(1)
}
