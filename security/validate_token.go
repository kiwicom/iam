package security

import (
	"errors"

	"github.com/iam/security/secrets"
)

var (
	errUnathorised = errors.New("incorrect token")
)

// VerifyToken accepts a token and a service struct and verifies if this token is accepted
func VerifyToken(secretManager secrets.SecretManager, service Service, requestToken string) error {
	if requestToken == "" {
		return errUnathorised
	}

	token, err := secretManager.GetAppToken(service.Name, service.Environment)
	if err != nil {
		return err
	}

	if token != requestToken {
		return errUnathorised
	}

	return nil
}
