package okta

import (
	gourl "net/url"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	cfg "github.com/iam/config"
	"github.com/iam/storage"
)

const iamGroupPrefix = "iam-"

// Group represents an Okta group
type Group struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
}

type oktaGroupProfile struct {
	Name        string
	Description string
}

func (c *Client) fetchGroups(userID, since string) ([]Group, error) {
	var filter string
	if since != "" {
		filter = "?filter=" + gourl.QueryEscape("lastMembershipUpdated gt \""+since+"\"")
	}

	var url string
	var err error
	if userID != "" {
		url, err = joinURL(c.baseURL, "/users/", userID, "/groups/")
	} else {
		url, err = joinURL(c.baseURL, "/groups/")
	}
	if err != nil {
		return nil, err
	}

	var allGroups []Group
	var resources []struct {
		ID                    string
		Profile               oktaGroupProfile
		LastMembershipUpdated time.Time
		LastFetched           time.Time
	}

	responses, err := c.fetchPagedResource(url + filter)
	if err != nil {
		return nil, err
	}

	for _, response := range responses {
		jsonErr := json.UnmarshalFromString(response, &resources)
		if jsonErr != nil {
			return nil, jsonErr
		}

		var groups []Group
		for i := range resources {
			group := &resources[i]
			if strings.HasPrefix(group.Profile.Name, iamGroupPrefix) {
				groups = append(groups, Group{
					ID:                    group.ID,
					Name:                  group.Profile.Name,
					Description:           group.Profile.Description,
					LastMembershipUpdated: group.LastMembershipUpdated,
				})
			}
		}
		allGroups = append(allGroups, groups...)
	}

	return allGroups, nil
}

func (c *Client) getUserGroups(user *User) ([]Group, error) {
	if user.GroupMembership != nil {
		// Group membership is cached, don't fetch it.
		return user.GroupMembership, nil
	}

	lockName := user.Email + ":groupMembership"
	// Deduplicate network calls and cache writes if this function is called
	// multiple times within the same instance.
	val, err, _ := c.group.Do(lockName, func() (interface{}, error) {
		lockErr := c.lock.Create(lockName)
		if lockErr == storage.ErrLockExists {
			// If there was a lock, it means another instance was fetching this data
			// recently, in that case, we should be able to just get the data from
			// cache.
			var u User
			if err := c.cache.Get(user.Email, &u); err != nil {
				return nil, err
			}
			return u.GroupMembership, nil
		}
		defer c.lock.Delete(lockName)

		groups, fetchErr := c.fetchGroups(user.OktaID, "")
		if fetchErr != nil {
			return nil, fetchErr
		}
		user.GroupMembership = groups

		cacheErr := c.cache.Set(user.Email, user, cfg.Expirations.User)
		if cacheErr != nil {
			raven.CaptureError(cacheErr, nil)
		}
		return groups, nil
	})

	if err != nil {
		return nil, err
	}
	return val.([]Group), nil
}
