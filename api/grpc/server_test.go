package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	pb "github.com/kiwicom/iam/api/grpc/v1"
	"github.com/kiwicom/iam/internal/services/okta"
)

type mockOktaService struct {
	mock.Mock
}

func (o *mockOktaService) AddPermissions(user *okta.User, service string) error {
	argsToReturn := o.Called(user, service)
	return argsToReturn.Error(0)
}

func (o *mockOktaService) GetUser(email string) (okta.User, error) {
	argsToReturn := o.Called(email)
	return argsToReturn.Get(0).(okta.User), argsToReturn.Error(1)
}

var testUser = okta.User{
	EmployeeNumber: "1",
	Email:          "test@test.com",
	FirstName:      "Test",
	LastName:       "Tester",
	Position:       "QA Tester",
	Permissions:    []string{"action:read"},
}

var wantUser = &pb.UserResponse{
	EmployeeNumber: 1,
	Email:          "test@test.com",
	FirstName:      "Test",
	LastName:       "Tester",
	Position:       "QA Tester",
	Permissions:    []string{"action:read"},
	Boocsek:        &pb.BoocsekAttributes{},
}

func TestHappyPath(t *testing.T) {
	userService := &mockOktaService{}
	server := &Server{userService: userService}

	userService.On("GetUser", "test@test.com").Once().Return(testUser, nil)
	userService.On("AddPermissions", &testUser, "service").Once().Return(nil)

	ctx := context.Background()
	md := metadata.New(map[string]string{"service-agent": "service/0 (Kiwi.com test)"})
	ctx = metadata.NewIncomingContext(ctx, md)

	gotUser, err := server.User(ctx, &pb.UserRequest{
		Email: "test@test.com",
	})

	assert.NoError(t, err, "shouldn't return error")
	assert.Equal(t, wantUser, gotUser, "Returns correct body")
	userService.AssertExpectations(t)
}
