package client

import (
	"fmt"
)

type TeamUser struct {
	ID       string `json:"id,omitempty"`
	UserName string `json:"userName,omitempty"`
	Email    string `json:"email,omitempty"`
}

// Team spec
type Team struct {
	ID      string     `json:"_id,omitempty"`
	Name    string     `json:"name,omitempty"`
	Type    string     `json:"type,omitempty"`
	Account string     `json:"account,omitempty"`
	Tags    []string   `json:"tags,omitempty"`
	Users   []TeamUser `json:"users,omitempty"`
}

// NewTeam Codefresh API expects a list of users IDs to create a new team
type NewTeam struct {
	ID      string   `json:"_id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	Account string   `json:"account,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Users   []string `json:"users,omitempty"`
}

// GetID implement CodefreshObject interface
func (team *Team) GetID() string {
	return team.ID
}

func (client *Client) GetTeamList() ([]Team, error) {
	fullPath := "/team"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var teams []Team

	err = DecodeResponseInto(resp, &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (client *Client) GetTeamByName(name string) (*Team, error) {

	teams, err := client.GetTeamList()
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		if team.Name == name {
			return &team, nil
		}
	}

	return nil, nil
}

func (client *Client) GetTeamByID(id string) (*Team, error) {

	teams, err := client.GetTeamList()
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		if team.ID == id {
			return &team, nil
		}
	}

	return nil, nil
}

//
func ConvertToNewTeam(team *Team) *NewTeam {
	var users []string

	for _, user := range team.Users {
		users = append(users, user.ID)
	}

	return &NewTeam{
		ID:      team.ID,
		Name:    team.Name,
		Type:    team.Type,
		Account: team.Account,
		Tags:    team.Tags,
		Users:   users,
	}
}

// CreateTeam POST team
func (client *Client) CreateTeam(team *Team) (*NewTeam, error) {

	newTeam := ConvertToNewTeam(team)
	body, err := EncodeToJSON(newTeam)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/team",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respTeam NewTeam
	err = DecodeResponseInto(resp, &respTeam)
	if err != nil {
		return nil, err
	}

	return &respTeam, nil
}

// DeleteTeam
func (client *Client) DeleteTeam(id string) error {
	fullPath := fmt.Sprintf("/team/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) SynchronizeClientWithGroup(name, ssoType string, notifications bool) error {

	fullPath := fmt.Sprintf("/team/group/synchronize/name/%s/type/%s?disableNotifications=%t", name, ssoType, notifications)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) AddUserToTeam(teamID, userID string) error {

	fullPath := fmt.Sprintf("/team/%s/%s/assignUserToTeam", teamID, userID)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteUserFromTeam(teamID, userID string) error {

	fullPath := fmt.Sprintf("/team/%s/%s/deleteUserFromTeam", teamID, userID)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) RenameTeam(teamID, name string) error {

	fullPath := fmt.Sprintf("/team/%s/renameTeam", teamID)

	team := Team{Name: name}
	body, err := EncodeToJSON(team)

	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
		Body:   body,
	}

	_, err = client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func GetUsersDiff(desiredUsers []string, existingUsers []TeamUser) (usersToAdd []string, usersToDelete []string) {

	existingUsersIDs := []string{}
	usersToAdd = []string{}
	usersToDelete = []string{}

	for _, user := range existingUsers {
		existingUsersIDs = append(existingUsersIDs, user.ID)
	}

	for _, id := range existingUsersIDs {
		ok := FindInSlice(desiredUsers, id)
		if !ok {
			usersToDelete = append(usersToDelete, id)
		}
	}

	for _, id := range desiredUsers {
		ok := FindInSlice(existingUsersIDs, id)
		if !ok {
			usersToAdd = append(usersToAdd, id)
		}
	}

	return usersToAdd, usersToDelete
}