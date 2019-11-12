package secrets

import (
	"errors"

	"github.com/spf13/viper"
)

// LocalSecretManager is just a struct
type LocalSecretManager struct{}

// CreateNewLocalSecretManager creates a new secret manager hooked up to Viper
func CreateNewLocalSecretManager() *LocalSecretManager {
	return &LocalSecretManager{}
}

// DoesTokenExist checks if a token is present in the secret manager
func (s LocalSecretManager) DoesTokenExist(reqToken string) bool {
	token := viper.GetString("TOKEN")
	return reqToken == token
}

// GetSetting gets a setting from Viper
func (s LocalSecretManager) GetSetting(key string) (string, error) {
	setting := viper.GetString(key)

	if setting == "" {
		return "", errors.New("setting not found")
	}

	return setting, nil
}
