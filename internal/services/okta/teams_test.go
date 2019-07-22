package okta

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCache struct {
	mock.Mock
	data map[string]interface{}
}

func (c *mockCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.data[key] = value
	return nil
}

func TestCacheTeams(t *testing.T) {
	cache := &mockCache{data: make(map[string]interface{})}
	users := []User{
		{
			FirstName:      "User1",
			TeamMembership: []string{"team1", "team2"},
		},
		{
			FirstName:      "User2",
			TeamMembership: []string{"team1", "team3"},
		},
	}

	expected := map[string]int{
		"team1": 2,
		"team2": 1,
		"team3": 1,
	}
	nTeams, err := cacheTeams(cache, users)
	assert.NoError(t, err)
	assert.Equal(t, 3, nTeams)
	assert.Equal(t, expected, cache.data["teams"])
}
