package grpc

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/kiwicom/iam/api/grpc/v1"
	"github.com/kiwicom/iam/internal/services/okta"
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
func (s *Server) User(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, userErr := s.userService.GetUser(in.Email)
	if userErr != nil {
		return nil, userErr
	}
	employeeNumber, intErr := strconv.ParseInt(user.EmployeeNumber, 10, 64)
	if intErr != nil {
		return nil, errors.New("unexpected server error")
	}

	attributes := pb.BoocsekAttributes{
		Site:        user.BoocsekAttributes.Site,
		Position:    user.BoocsekAttributes.Position,
		Channel:     user.BoocsekAttributes.Channel,
		Tier:        user.BoocsekAttributes.Tier,
		Team:        user.BoocsekAttributes.Team,
		TeamManager: user.BoocsekAttributes.TeamManager,
		Staff:       user.BoocsekAttributes.Staff,
		State:       user.BoocsekAttributes.State,
		KiwibaseId:  user.BoocsekAttributes.KiwibaseID,
		Substate:    user.BoocsekAttributes.Substate,
		Skills:      user.BoocsekAttributes.Skills,
	}

	return &pb.UserResponse{
		EmployeeNumber: employeeNumber,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Position:       user.Position,
		Department:     user.Department,
		Location:       user.Location,
		Manager:        user.Manager,
		TeamMembership: user.TeamMembership,
		Boocsek:        &attributes,
	}, nil
}
