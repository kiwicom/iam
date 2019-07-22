package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	tests := map[string]struct {
		token  string
		errors bool
	}{
		"Bearer token": {
			token:  "token",
			errors: false,
		},
		"token": {
			token:  "",
			errors: true,
		},
		"Bearer Bearer token": {
			token:  "Bearer token",
			errors: false,
		},
		"Bearer": {
			token:  "",
			errors: false,
		},
		"unexpected Bearer token": {
			token:  "",
			errors: true,
		},
	}

	for test, result := range tests {
		actual, err := GetToken(test)
		assert.Equal(t, result.token, actual)

		if result.errors {
			assert.Error(t, err)
		}
	}
}
