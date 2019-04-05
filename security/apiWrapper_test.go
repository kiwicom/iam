package security

import (
	"errors"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedSecretManager struct {
	mock.Mock
}

func (s *mockedSecretManager) GetAppToken(app string) (string, error) {
	if app == "serviceName" {
		return "valid token", nil
	}
	return "", errors.New("wrong token bro")
}

func (s *mockedSecretManager) GetSetting(app string) (string, error) {
	if app == "serviceName" {
		return "valid token", nil
	}
	return "", errors.New("wrong token bro")
}

func createFakeManager() SecretManager {
	return &mockedSecretManager{}
}

func TestCheckAuth(t *testing.T) {
	secrets := createFakeManager()

	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	assert.Error(t, checkAuth(req, secrets), "Should error on missing email")

	req, _ = http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)
	assert.Error(t, checkAuth(req, secrets), "Should error on missing User-Agent")
	req.Header.Set("User-Agent", "serviceName")

	assert.Error(t, checkAuth(req, secrets), "Should error on missing Authorization header")
	req.Header.Set("Authorization", "invalid token")

	assert.Error(t, checkAuth(req, secrets), "Should error on invalid token")
	req.Header.Set("Authorization", "valid token")
	viper.Set("TOKEN_serviceName_OKTA", "valid token")

	assert.NoError(t, checkAuth(req, secrets), "Should not error on valid request token")
}
