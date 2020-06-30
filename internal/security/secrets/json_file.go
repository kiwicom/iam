package secrets

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"	
)

// Secrets represents the JSON file structure 
type Secrets struct {
	Settings map[string]string            `json:"settings"`
	TokenMap map[string]map[string]string `json:"tokens"`
}

// JSONFileManager holds a local copy of all secrets (settings & S2S tokens)
type JSONFileManager struct {
	raw []byte
	settings map[string]string
	tokens map[string]bool
}

// CreateNewJSONFileManager creates a new secret manager hooked up to Viper
func CreateNewJSONFileManager(path string) (*JSONFileManager, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	log.Println("Using JSON secret file at:", path)

	return &JSONFileManager{
		raw: data,
	}, nil
}

// SyncSecrets syncs all the available tokens from env and saves them to local state
func (s *JSONFileManager) SyncSecrets() error {
	var secrets Secrets

	if err := json.Unmarshal(s.raw, &secrets); err != nil {
		return err
	}

	s.settings = secrets.Settings

	mappedTokens := make(map[string]bool)
	tokenCount := 0

	for _, app := range secrets.TokenMap {
		for _, token := range app {
			tokenCount++
			mappedTokens[token] = true
		}
	}
	s.tokens = mappedTokens

	log.Printf("Synced %v settings and %v tokens (%v unique).", len(s.settings), tokenCount, len(s.tokens))

	return nil
}

// DoesTokenExist checks if a token is present in the secret manager
func (s JSONFileManager) DoesTokenExist(reqToken string) bool {
	return s.tokens[reqToken]
}

// GetSetting gets a setting from the secret manager
func (s JSONFileManager) GetSetting(key string) (string, error) {
	data := s.settings[key]

	if data == "" {
		return "", errors.New("key '" + key + "' not found in SecretManager")
	}

	return data, nil
}
