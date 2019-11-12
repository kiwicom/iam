package security

import (
	"errors"

	"github.com/kiwicom/iam/internal/security/secrets"
)

var (
	errUnathorised = errors.New("incorrect token")
)

// VerifyToken accepts a token and a service struct and verifies if this token is accepted
func VerifyToken(secretManager secrets.SecretManager, service Service, requestToken string) error {
	if requestToken == "" {
		return errUnathorised
	}

	exists := secretManager.DoesTokenExist(requestToken)
	if !exists {
		return errUnathorised
	}

	return nil
}
