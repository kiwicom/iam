package okta

import (
	"errors"
	"time"

	"github.com/getsentry/raven-go"

	"github.com/kiwicom/iam/internal/storage"
)

// Permissions contains a map with names as keys and lists of IDs as values.
type Permissions map[string][]string

// ErrNotReady is returned when data is not ready and the client should retry
// the request later
var ErrNotReady = errors.New("data not ready")

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

		// This case can happen if data is lost between syncs, if a sync fails, or
		// if the cache version is bumped and a request is received before syncing
		// groups. We try syncing groups from Okta again, and delete the
		// sync-timestamp before doing so. This is to ensure we retrieve all the
		// data. An error is returned to the client, because if we would keep the
		// connection open while syncing groups, it would timeout and fail anyway.
		go func() {
			_, _, _ = c.group.Do("resync-groups", func() (interface{}, error) {
				err := c.cache.Del("groups-sync-timestamp")
				if err != nil {
					raven.CaptureError(err, nil)
				}
				c.SyncGroups()
				return nil, nil
			})
		}()

		return nil, ErrNotReady
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

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// GetUserPermissions returns permissions for the specified user for the
// requested services.
func (c *Client) GetUserPermissions(email string, services []string) (map[string][]string, error) {
	allPermissions, err := c.GetServicesPermissions(services)
	if err != nil {
		return nil, err
	}

	userPermissions := make(map[string][]string)
	for _, service := range services {
		for permission, users := range allPermissions[service] {
			if stringInSlice(email, users) {
				userPermissions[service] = append(userPermissions[service], permission)
			}
		}
	}

	return userPermissions, nil
}
