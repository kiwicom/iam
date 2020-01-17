package cfg

import (
	"time"
)

// Expirations contains expiration times for anything that needs to be cached
var Expirations = struct {
	User             time.Duration
	GroupMemberships time.Duration
	GroupsLastSync   time.Duration
}{
	User:             time.Hour * 24,
	GroupMemberships: time.Hour * 24,
	// GroupsLastSync does not expire, but its value is deleted once a day (and on
	// new deploys) in order to refetch all groups. This cannot be done through
	// setting an expiration because GroupsLastSync is set every 10 minutes, and
	// if there is a value for it, it's used to fetch only the changes since the
	// last sync instead of all groups.
	GroupsLastSync: 0,
}
