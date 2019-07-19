package okta

import (
	"errors"
	"regexp"
	"strings"

	"github.com/getsentry/raven-go"
	cfg "github.com/iam/config"
	"github.com/iam/storage"
)

// GroupMembership holds the current user ids for users who are part of a given group
type GroupMembership struct {
	GroupID   string
	GroupName string
	Users     []string
}

func (c *Client) fetchGroupMembership(groupID string) ([]string, error) {
	url, err := joinURL(c.baseURL, "/groups/", groupID, "/users")
	if err != nil {
		return nil, err
	}

	var allUsers []string
	var resources []struct {
		Profile oktaUserProfile
	}

	responses, err := c.fetchPagedResource(url)
	if err != nil {
		return nil, err
	}

	for _, response := range responses {
		jsonErr := json.UnmarshalFromString(response, &resources)
		if jsonErr != nil {
			return nil, jsonErr
		}

		users := make([]string, len(resources))
		for i := range resources {
			users[i] = resources[i].Profile.EmployeeNumber
		}

		allUsers = append(allUsers, users...)
	}

	return allUsers, nil
}

func (c *Client) fetchGroupMemberships(groups []Group) ([]GroupMembership, error) {
	groupMemberships := make([]GroupMembership, len(groups))

	for i, group := range groups {
		users, fetchErr := c.fetchGroupMembership(group.ID)
		if fetchErr != nil {
			return nil, fetchErr
		}
		groupMemberships[i] = GroupMembership{
			group.ID,
			group.Name,
			users,
		}
	}

	return groupMemberships, nil
}

var groupPattern = regexp.MustCompile(`^iam-[\w-]+\.([\w-]+\.?)+$`)

func (c *Client) updateGroupMemberships(memberships []GroupMembership) error {
	cachedGroupMemberships := make(map[string]map[string]bool)

	for _, membership := range memberships {
		if !groupPattern.Match([]byte(membership.GroupName)) {
			formatErr := errors.New("group name has incorrect format: " + membership.GroupName)
			raven.CaptureError(formatErr, nil)
			return formatErr
		}

		// iam-serviceName:rule
		groupParts := strings.SplitAfterN(membership.GroupName, ".", 2)
		serviceName := groupMembershipPrefix + strings.Replace(strings.TrimRight(groupParts[0], "."), iamGroupPrefix, "", 1)

		err := c.cache.Get(serviceName, &cachedGroupMemberships)
		if err != nil {
			if err != storage.ErrNotFound {
				return err
			}
		}

		if _, ok := cachedGroupMemberships[groupParts[1]]; !ok {
			cachedGroupMemberships[groupParts[1]] = make(map[string]bool)
		}

		for _, userid := range membership.Users {
			cachedGroupMemberships[groupParts[1]][userid] = true
		}

		if err := c.cache.Set(serviceName, cachedGroupMemberships, cfg.Expirations.GroupMemberships); err != nil {
			return err
		}
	}

	return nil
}
