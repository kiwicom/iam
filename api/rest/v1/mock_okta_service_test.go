package rest

import (
	"github.com/stretchr/testify/mock"

	"github.com/kiwicom/iam/internal/services/okta"
)

type mockOktaService struct {
	mock.Mock
}

func (o *mockOktaService) AddPermissions(user *okta.User, service string) error {
	argsToReturn := o.Called(user, service)
	return argsToReturn.Error(0)
}

func (o *mockOktaService) GetTeams() (map[string]int, error) {
	args := o.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (o *mockOktaService) GetUser(email string) (okta.User, error) {
	argsToReturn := o.Called(email)
	return argsToReturn.Get(0).(okta.User), argsToReturn.Error(1)
}

func (o *mockOktaService) GetServicesPermissions(services []string) (map[string]okta.Permissions, error) {
	argsToReturn := o.Called(services)
	return argsToReturn.Get(0).(map[string]okta.Permissions), argsToReturn.Error(1)
}

func (o *mockOktaService) GetUserPermissions(email string, services []string) (map[string][]string, error) {
	argsToReturn := o.Called(email, services)
	return argsToReturn.Get(0).(map[string][]string), argsToReturn.Error(1)
}

func (o *mockOktaService) GetGroups() ([]okta.Group, error) {
	argsToReturn := o.Called()
	return argsToReturn.Get(0).([]okta.Group), argsToReturn.Error(1)
}
