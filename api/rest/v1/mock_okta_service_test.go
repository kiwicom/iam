package rest

import (
	"github.com/kiwicom/iam/internal/services/okta"
	"github.com/stretchr/testify/mock"
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
