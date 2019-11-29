package okta

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiwicom/iam/internal/storage"
)

func TestGroupMemberships(t *testing.T) {
	client := NewClient(&ClientOpts{
		Cache: storage.NewInMemoryCache(),
	})

	tests := []struct {
		groupName   string
		shouldError bool
	}{
		{"iam-service.valid.permission", false},
		{"iam-long-service-name.with-long.permission", false},
		{"iam-service:invalid:permission", false},
		{"service", false},
	}

	for _, test := range tests {
		err := client.updateGroupMemberships([]GroupMembership{
			{"group-id", test.groupName, []string{"user1", "user2"}},
		})

		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGroupMembershipsWithInvalidNames(t *testing.T) {
	client := NewClient(&ClientOpts{
		Cache: storage.NewInMemoryCache(),
	})
	err := client.updateGroupMemberships([]GroupMembership{
		{"group-id", "service", []string{"user1", "user2"}},
		{"group-id", "iam-service.valid.permission", []string{"user1", "user2"}},
	})

	assert.NoError(t, err)
	permissions, err := client.GetServicePermissions("service")
	assert.NoError(t, err)
	assert.Contains(t, permissions["valid.permission"], "user1")
	assert.Contains(t, permissions["valid.permission"], "user2")
}

func TestGroupMembershipsInvalidation(t *testing.T) {
	cache := storage.NewInMemoryCache()
	client := NewClient(&ClientOpts{Cache: cache})

	err := client.updateGroupMemberships([]GroupMembership{
		{"group-id", "iam-service.permission1", []string{"user1", "user2"}},
		{"group-id", "iam-service.permission2", []string{"user1", "user2"}},
	})

	membershipsBefore := make(map[string]map[string]bool)
	_ = cache.Get("group-membership:service", &membershipsBefore)

	assert.NoError(t, err)
	assert.Equal(t, map[string]map[string]bool{
		"permission1": {
			"user1": true,
			"user2": true,
		},
		"permission2": {
			"user1": true,
			"user2": true,
		},
	}, membershipsBefore, "Group memberships are added correctly")

	err = client.updateGroupMemberships([]GroupMembership{
		{"group-id", "iam-service.permission1", []string{"user2"}},
		{"group-id", "iam-service.permission2", []string{"user1", "user2"}},
	})

	membershipsAfter := make(map[string]map[string]bool)
	_ = cache.Get("group-membership:service", &membershipsAfter)

	assert.NoError(t, err)
	assert.Equal(t, map[string]map[string]bool{
		"permission1": {
			"user2": true,
		},
		"permission2": {
			"user1": true,
			"user2": true,
		},
	}, membershipsAfter, "Group membership is invalidated correctly")
}
