package security

import (
	"errors"
	"strings"
)

// GetToken extracts token from Bearer scheme
func GetToken(authorization string) (string, error) {
	token := strings.SplitN(authorization, "Bearer", 2)

	if len(token) == 2 && token[0] == "" {
		return strings.TrimSpace(token[1]), nil
	}

	return "", errors.New("invalid auth scheme")
}
