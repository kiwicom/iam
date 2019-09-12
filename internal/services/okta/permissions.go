package okta

import (
	"time"

	"github.com/getsentry/raven-go"

	"github.com/kiwicom/iam/internal/storage"
)

// Permissions contains a map with names as keys and lists of IDs as values.
type Permissions map[string][]string

// GetServicePermissions retrieves permissions for the specified service.
func (c *Client) GetServicePermissions(service string) (Permissions, error) {
	permissions := make(Permissions)
	cachedGroupMemberships := make(map[string]map[string]bool)

	err := c.cache.Get(groupMembershipPrefix+service, &cachedGroupMemberships)
	if err != nil {
		if err != storage.ErrNotFound {
			return permissions, err
		}

		timestamp := time.Time{}
		_ = c.cache.Get("groups-sync-timestamp", &timestamp)
		if time.Now().Before(timestamp.Add(10 * time.Minute)) {
			// If there are no groups cached for the service and it's less than 10
			// minutes from the last sync, we assume that there are no groups for that
			// service.
			return permissions, nil
		}

		// This case can happen if data is lost between syncs, or if a sync fails.
		// We try syncing groups from Okta again, and delete the sync-timestamp
		// before doing so. This is to ensure the synced data is correct.
		// If a request ends up here it will in most (if not all) cases timeout,
		// not sure if a better solution would be to return nothing or an error.
		val, err, _ := c.group.Do("resync-groups", func() (interface{}, error) {
			err := c.cache.Del("groups-sync-timestamp")
			if err != nil {
				raven.CaptureError(err, nil)
			}
			c.SyncGroups()
			return c.GetServicePermissions(service)
		})

		if err != nil {
			return permissions, err
		}
		return val.(Permissions), nil
	}

	for groupName, users := range cachedGroupMemberships {
		permissions[groupName] = make([]string, 0, len(users))
		for userID := range users {
			permissions[groupName] = append(permissions[groupName], userID)
		}
	}

	return permissions, nil
}

// GetServicesPermissions returns all permissions for the requested services, for
// each permissions there is a list of IDs that have such permissions.
func (c *Client) GetServicesPermissions(services []string) (map[string]Permissions, error) {
	permissions := make(map[string]Permissions)

	for _, service := range services {
		p, err := c.GetServicePermissions(service)
		if err != nil {
			return nil, err
		}
		permissions[service] = p
	}

	return permissions, nil
}
