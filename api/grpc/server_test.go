package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/kiwicom/iam/api"
	pb "github.com/kiwicom/iam/api/grpc/v1"
	"github.com/kiwicom/iam/internal/services/okta"
)

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
	userService := &api.MockOktaService{}
	request, _ := http.NewRequest("GET", "/?email=test@test.com&permissions=true", nil)
	request.Header.Set("User-Agent", "service/0 (Kiwi.com test)")
	server := &Server{userService: userService}

	userService.On("GetUser", "test@test.com").Once().Return(testUser, nil)
	userService.On("AddPermissions", &testUser, "service").Once().Return(nil)

	ctx := context.Background()
	md := metadata.New(map[string]string{"service-agent": "service"})
	ctx = metadata.NewIncomingContext(ctx, md)

	gotUser, err := server.User(ctx, &pb.UserRequest{
		Email: "test@test.com",
	})

	assert.NoError(t, err, "shouldn't return error")
	assert.Equal(t, wantUser, gotUser, "Returns correct body")
	userService.AssertExpectations(t)
}
