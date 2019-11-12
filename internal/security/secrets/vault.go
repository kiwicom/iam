package secrets

import (
	"errors"
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
)

type localStorage struct {
	// Tokens for apps integrating with IAM
	tokens map[string]bool

	// Settings for IAM
	settings map[string]string
}

// VaultManager is the local struct for Vault connection
type VaultManager struct {
	client    *vault.Client
	storage   localStorage
	namespace string
}

const requiredPrefix = "/secret"

// CreateNewVaultClient create a new client to connect to Vault
func CreateNewVaultClient(address, token, namespace string) (*VaultManager, error) {
	if address == "" {
		return nil, errors.New("missing Vault address")
	}

	if token == "" {
		return nil, errors.New("missing Vault token")
	}

	if namespace == "" {
		return nil, errors.New("missing Vault namespace setting")
	}

	client, _ := vault.NewClient(&vault.Config{
		Address: address,
	})

	client.SetToken(token)

	client.SetNamespace(namespace)

	return &VaultManager{client: client, namespace: namespace}, nil
}

// SyncAppTokens syncs all the available tokens from Vault and saves them to local state
func (s *VaultManager) SyncAppTokens() error {
	data, err := s.fetchData("app_tokens")

	if err != nil {
		return err
	}

	mappedTokens := make(map[string]bool, len(data))

	for _, value := range data {
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid conversion to string for value <%v> of type <%T>", v, v)
		}
		mappedTokens[v] = true
	}

	s.storage.tokens = mappedTokens

	return nil
}

// DoesTokenExist checks if a token is present in the secret manager
func (s *VaultManager) DoesTokenExist(token string) bool {
	return s.storage.tokens[token]
}

// SyncAppSettings syncs all application settings from Vault and saves them locally
func (s *VaultManager) SyncAppSettings() error {
	data, err := s.fetchData("settings")

	if err != nil {
		log.Println(err)
		return err
	}

	mappedSettings := make(map[string]string, len(data))

	for key, value := range data {
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid conversion to string for value <%v> of type <%T>", v, v)
		}
		mappedSettings[key] = v
	}

	s.storage.settings = mappedSettings

	return nil
}

// GetSetting returns a single setting value
func (s *VaultManager) GetSetting(key string) (string, error) {
	data := s.storage.settings[key]

	if data == "" {
		return "", errors.New("key '" + key + "' not found in SecretManager")
	}

	return data, nil
}

func (s *VaultManager) fetchData(subpath string) (map[string]interface{}, error) {
	response, err := s.client.Logical().Read(requiredPrefix + "/" + s.namespace + "/" + subpath)

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("empty response from Vault")
	}

	return response.Data, nil
}
