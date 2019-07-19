package grpc

import (
	"context"
	"log"

	"github.com/iam/security"
	"github.com/iam/security/secrets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	errBadUA           = status.Errorf(codes.Unauthenticated, "invalid service-agent")
)

// Metadata keys for the headers/trailers.
const (
	metadataUserAgent     = "service-agent"
	metadataAuthorization = "authorization"
)

// UnarySecurityWrapper creates a new Security middleware for gRPC. It will check for the presence of a useragent.
// It will also validate that the sent token is correct.
func UnarySecurityWrapper(secretManager secrets.SecretManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		// service-agent is used as gRPC tools currently don't allow for overriding user-agent
		if len(md[metadataUserAgent]) == 0 {
			return nil, errBadUA
		}

		service, serviceErr := security.GetService(md[metadataUserAgent][0])

		if serviceErr != nil {
			return nil, errBadUA
		}

		if len(md[metadataAuthorization]) == 0 {
			return nil, errInvalidToken
		}

		token, err := security.GetToken(md[metadataAuthorization][0])
		if err != nil {
			return nil, errInvalidToken
		}

		tokenErr := security.VerifyToken(secretManager, service, token)
		if tokenErr != nil {
			log.Println(tokenErr)
			return nil, errInvalidToken
		}

		m, err := handler(ctx, req)
		if err != nil {
			log.Printf("RPC failed with error %v", err)
		}
		return m, err
	}
}
