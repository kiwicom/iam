package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kiwicom/iam/internal/security/secrets"
)

type mockedSecretManager struct {
	mock.Mock
}

func (s *mockedSecretManager) GetAppToken(app, environment string) (string, error) {
	if app == "serviceName" && environment == "environment" {
		return "valid token", nil
	}
	return "", errors.New("wrong token bro")
}

func (s *mockedSecretManager) GetSetting(_ string) (string, error) {
	return "", nil
}

func createFakeManager() secrets.SecretManager {
	return &mockedSecretManager{}
}

type mockedMetricsService struct {
	mock.Mock
}

func (m *mockedMetricsService) Incr(serviceName string, tags ...string) {
	m.Called(serviceName, tags)
}

func TestUnhappyPathCheckAuth(t *testing.T) {
	m := &mockedMetricsService{}
	sm := createFakeManager()
	s := Server{
		SecretManager: sm,
		MetricClient:  m,
	}

	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	err := s.checkAuth(req)
	assert.Error(t, err, "Should error on missing email")

	req, _ = http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)
	err = s.checkAuth(req)
	assert.Error(t, err, "Should error on missing User-Agent")

	req.Header.Set("User-Agent", "serviceName/version (Kiwi.com environment)")
	err = s.checkAuth(req)
	assert.Error(t, err, "Should error on missing Authorization header")

	req.Header.Set("Authorization", "invalid token")
	err = s.checkAuth(req)
	assert.Error(t, err, "Should error on invalid token schema")
	m.AssertNotCalled(t, "Incr")

	req.Header.Set("Authorization", "Bearer invalid token")
	err = s.checkAuth(req)
	assert.Error(t, err, "Should error on invalid token")
	m.AssertNotCalled(t, "Incr")
}

func TestHappyPathCheckAuth(t *testing.T) {
	m := &mockedMetricsService{}
	sm := createFakeManager()
	s := Server{
		SecretManager: sm,
		MetricClient:  m,
	}

	req, _ := http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)

	req.Header.Set("User-Agent", "serviceName/version (Kiwi.com environment)")
	err := s.checkAuth(req)
	assert.Error(t, err, "Should error on missing Authorization header")
	m.AssertNotCalled(t, "Incr")

	req.Header.Set("Authorization", "Bearer valid token")
	m.On("Incr", "incoming.requests", []string{"service-name:servicename", "service-environment:environment"})
	err = s.checkAuth(req)
	assert.NoError(t, err, "Should not error on valid request token")
	m.AssertNumberOfCalls(t, "Incr", 1)
}
