package secrets

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// LocalSecretManager is just a struct
type LocalSecretManager struct{}

// CreateNewLocalSecretManager creates a new secret manager hooked up to Viper
func CreateNewLocalSecretManager() *LocalSecretManager {
	return &LocalSecretManager{}
}

// GetAppToken gets the token from Viper
func (s LocalSecretManager) GetAppToken(app, environment string) (string, error) {
	tokenName := strings.ToUpper(app) + "_" + strings.ToUpper(environment)
	token := viper.GetString("TOKEN_" + tokenName)

	if token == "" {
		return "", errors.New("token '" + tokenName + "' not found")
	}

	return token, nil
}

// GetSetting gets a setting from Viper
func (s LocalSecretManager) GetSetting(key string) (string, error) {
	setting := viper.GetString(key)

	if setting == "" {
		return "", errors.New("setting not found")
	}

	return setting, nil
}
