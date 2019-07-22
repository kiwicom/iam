package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, stringInSlice("some", []string{"where", "some"}))
	assert.True(t, stringInSlice("caseinsensitive", []string{"CaseInsensitive"}))
	assert.False(t, stringInSlice("some", []string{"where", "somehere"}))
}
