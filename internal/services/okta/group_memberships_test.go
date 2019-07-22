package okta

import (
	"testing"

	"github.com/iam/internal/storage"
	"github.com/stretchr/testify/assert"
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
		{"iam-service:invalid:permission", true},
		{"service", true},
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
