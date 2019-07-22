package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetService(t *testing.T) {
	tests := map[string][]string{
		"BALKAN/4704b82 (Kiwi.com sandbox)":    {"BALKAN", "sandbox"},
		"balkan/1.42.1 (Kiwi.com sandbox)":     {"balkan", "sandbox"},
		"balkan-graphql/1.42.1 (Kiwi.com dev)": {"balkan-graphql", "dev"},
		"balkan_graphql/1.42.1 (Kiwi.com dev)": {"balkan_graphql", "dev"},
	}

	for test, expected := range tests {
		res, err := GetService(test)
		assert.Equal(t, expected[0], res.Name)
		assert.Equal(t, expected[1], res.Environment)
		assert.NoError(t, err)
	}

	invalidUAs := []string{
		"balkan",
		"balkan graphql/1.42.1 (Kiwi.com test)",
		"",
	}

	for _, ua := range invalidUAs {
		res, err := GetService(ua)
		assert.Equal(t, "", res.Name)
		assert.Equal(t, "", res.Environment)
		assert.Error(t, err)
	}
}

func TestCheckServiceName(t *testing.T) {
	tests := map[string]bool{
		"balkan":            false,
		"balkan_PROD1-test": false,
		"balkan%2f../":      true,
		"balkan/../":        true,
		"":                  true,
		"balkan$":           true,
	}

	for input, shouldError := range tests {
		if shouldError {
			assert.Error(t, CheckServiceName(input))
		} else {
			assert.NoError(t, CheckServiceName(input))
		}
	}
}
