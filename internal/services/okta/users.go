package okta

import (
	"errors"
	"log"
	"net/http"
)

type oktaUserProfile struct {
	EmployeeNumber     string
	FirstName          string
	LastName           string
	Department         string
	Email              string
	UserType           string
	SfOrgStructure     string `json:"SF_orgStructure"`
	SfJobTitle         string `json:"SF_jobTitle"`
	SfLocation         string `json:"SF_location"`
	Manager            string
	BoocsekSite        string   `json:"boocsek_site"`
	BoocsekPosition    string   `json:"boocsek_position"`
	BoocsekChannel     string   `json:"boocsek_channel"`
	BoocsekTier        string   `json:"boocsek_tier"`
	BoocsekTeam        string   `json:"boocsek_team"`
	BoocsekTeamManager string   `json:"boocsek_team_manager"`
	BoocsekStaff       string   `json:"boocsek_staff"`
	BoocsekState       string   `json:"boocsek_state"`
	BoocsekKiwibaseID  int32    `json:"boocsek_kiwibase_id"`
	BoocsekSubstate    string   `json:"boocsek_substate"`
	BoocsekSkills      []string `json:"boocsek_skills"`
}

func formatUser(oktaID string, user *oktaUserProfile) User {
	teamMembership := append(make([]string, 0), user.SfOrgStructure) // Deprecated
	skills := append(make([]string, 0), user.BoocsekSkills...)

	boocsekAttributes := BoocsekAttributes{
		Site:        user.BoocsekSite,
		Position:    user.BoocsekPosition,
		Channel:     user.BoocsekChannel,
		Tier:        user.BoocsekTier,
		Team:        user.BoocsekTeam,
		TeamManager: user.BoocsekTeamManager,
		Staff:       user.BoocsekStaff,
		State:       user.BoocsekState,
		KiwibaseID:  user.BoocsekKiwibaseID,
		Substate:    user.BoocsekSubstate,
		Skills:      skills,
	}

	return User{
		OktaID:                oktaID,
		EmployeeNumber:        user.EmployeeNumber,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Department:            user.Department,
		Email:                 user.Email,
		Position:              user.SfJobTitle,
		Location:              user.SfLocation,
		IsVendor:              user.UserType != "Regular Employee",
		TeamMembership:        teamMembership,
		OrganizationStructure: user.SfOrgStructure,
		Manager:               user.Manager,
		BoocsekAttributes:     boocsekAttributes,
	}
}

// ErrUserNotFound is returned when a user is not present in Okta
var ErrUserNotFound = errors.New("user not found")

// fetchUser retrieves a user from Okta by email
func (c *Client) fetchUser(email string) (User, error) {
	userURL, err := joinURL(c.baseURL, "/users/", email)
	if err != nil {
		return User{}, err
	}

	var response struct {
		ID      string
		Profile oktaUserProfile
	}
	var request = Request{
		Method: "GET",
		URL:    userURL,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := c.fetch(request)
	if err != nil {
		return User{}, err
	}
	if httpResponse.StatusCode == http.StatusNotFound {
		return User{}, ErrUserNotFound
	}
	if httpResponse.StatusCode != http.StatusOK {
		var errorMessage = "GET " + userURL + " returned error: " + httpResponse.Status
		log.Println(errorMessage)
		return User{}, errors.New(errorMessage)
	}

	jsonErr := httpResponse.JSON(&response)
	if jsonErr != nil {
		return User{}, jsonErr
	}

	var user = formatUser(response.ID, &response.Profile)
	return user, nil
}

// fetchAllUsers retrieves all Okta users
func (c *Client) fetchAllUsers() ([]User, error) {
	var allUsers []User

	url, err := joinURL(c.baseURL, "/users/")
	if err != nil {
		return nil, err
	}

	responses, err := c.fetchPagedResource(url)
	if err != nil {
		return nil, err
	}

	for _, response := range responses {
		var resources []struct {
			ID      string
			Profile oktaUserProfile
		}

		jsonErr := json.UnmarshalFromString(response, &resources)
		if jsonErr != nil {
			return nil, jsonErr
		}

		users := make([]User, len(resources))
		for i := range resources {
			user := &resources[i]
			users[i] = formatUser(user.ID, &user.Profile)
		}

		allUsers = append(allUsers, users...)
	}

	return allUsers, err
}
