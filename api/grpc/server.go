package grpc

import (
	"context"
	"strconv"

	pb "gitlab.skypicker.com/platform/security/iam/api/grpc/v1"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

type userDataService interface {
	GetUser(string) (okta.User, error)
	AddPermissions(*okta.User, string) error
}

// Server is an instance of the GRPC server struct which includes all dependencies
type Server struct {
	userService userDataService
}

// CreateServer creates a new Server struct and assigns all dependencies to it
func CreateServer(userServiceClient userDataService) (*Server, error) {
	return &Server{userService: userServiceClient}, nil
}

// User returns a single user based on email
func (s *Server) User(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, userErr := s.userService.GetUser(in.Email)
	if userErr != nil {
		return &pb.UserResponse{}, userErr
	}
	employeNumber, _ := strconv.ParseInt(user.EmployeeNumber, 10, 64)

	return &pb.UserResponse{
		EmployeeNumber: employeNumber,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Position:       user.Position,
		Department:     user.Department,
		Location:       user.Location,
		Manager:        user.Manager,
		TeamMembership: user.TeamMembership,
	}, nil
}