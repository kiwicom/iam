package okta

import (
	"time"

	cfg "github.com/iam/config"
)

type cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
}

// CacheTeams extracts teams from user profiles and saves them to cache
func cacheTeams(cache cache, users []User) (int, error) {
	teams := make(map[string]int)
	for i := 0; i < len(users); i++ {
		for _, team := range users[i].TeamMembership {
			teams[team]++
		}
	}

	err := cache.Set("teams", teams, cfg.Expirations.Teams)
	return len(teams), err
}
