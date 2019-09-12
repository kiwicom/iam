package okta

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kiwicom/iam/internal/storage"
)

func TestGetServicePermissions(t *testing.T) {
	client := NewClient(&ClientOpts{Cache: storage.NewInMemoryCache()})

	cachedGroups := map[string]map[string]bool{
		"access": {
			"user1": true,
			"user2": true,
		},
	}
	actualEmptyGroups := map[string]map[string]bool{}

	_ = client.cache.Set("groups-sync-timestamp", time.Now(), 0)
	_ = client.cache.Set(groupMembershipPrefix+"cached", cachedGroups, 0)
	_ = client.cache.Set(groupMembershipPrefix+"actual-empty", actualEmptyGroups, 0)

	permissions, err := client.GetServicePermissions("cached")
	assert.NoError(t, err)
	sort.Strings(permissions["access"])
	assert.Equal(t, Permissions{"access": {"user1", "user2"}}, permissions)

	permissions, err = client.GetServicePermissions("actual-empty")
	assert.NoError(t, err)
	assert.Equal(t, Permissions{}, permissions)

	permissions, err = client.GetServicePermissions("assumed-empty")
	assert.NoError(t, err)
	assert.Equal(t, Permissions{}, permissions)
}
