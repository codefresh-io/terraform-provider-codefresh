package cfclient

import (
	"fmt"
	"strings"
)

type Credentials struct {
	Permissions []string `json:"permissions,omitempty"`
}

type Login struct {
	Credentials Credentials `json:"credentials,omitempty"`
	PersonalGit bool        `json:"personalGit,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
	IDP         IDP         `json:"idp,omitempty"`
	Idp_ID      string      `json:"idp_id,omitempty"`
	Sso         bool        `json:"sso,omitempty"`
}

type ShortProfile struct {
	UserName string `json:"userName,omitempty"`
}

type PublicProfile struct {
	HasPassword bool `json:"hasPassword,omitempty"`
}

type Personal struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Country     string `json:"country,omitempty"`
}

type User struct {
	ID             string              `json:"_id,omitempty"`
	UserName       string              `json:"userName"`
	Email          string              `json:"email"`
	Personal       *Personal           `json:"personal,omitempty"`
	Roles          []string            `json:"roles,omitempty"`
	DefaultAccount int                 `json:"defaultAccount,omitempty"`
	Account        []Account           `json:"account,omitempty"`
	Status         string              `json:"status,omitempty"`
	RegisterDate   string              `json:"register_date,omitempty"`
	HasPassword    bool                `json:"hasPassword,omitempty"`
	Notifications  []NotificationEvent `json:"notifications,omitempty"`
	ShortProfile   ShortProfile        `json:"shortProfile,omitempty"`
	PublicProfile  PublicProfile       `json:"publicProfile,omitempty"`
	Logins         []Login             `json:"logins,omitempty"`
	InviteURL      string              `json:"inviteUrl,omitempty"`
}

type NewUser struct {
	ID       string    `json:"_id,omitempty"`
	UserName string    `json:"userName"`
	Email    string    `json:"email"`
	Logins   []Login   `json:"logins,omitempty"`
	Roles    []string  `json:"roles,omitempty"`
	Account  []string  `json:"account,omitempty"`
	Personal *Personal `json:"personal,omitempty"`
}

type UserAccounts struct {
	UserName string    `json:"userName"`
	Account  []Account `json:"account"`
}

// The API accepts two different schemas when updating the user details
func generateUserDetailsBody(userName, userEmail string) string {
	userDetails := fmt.Sprintf(`{"userDetails": "%s"}`, userEmail)
	if userName != "" {
		userDetails = fmt.Sprintf(`{"userName": "%s", "email": "%s"}`, userName, userEmail)
	}
	return userDetails
}

func (client *Client) AddNewUserToAccount(accountId, userName, userEmail string) (*User, error) {

	fullPath := fmt.Sprintf("/accounts/%s/adduser", accountId)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   []byte(generateUserDetailsBody(userName, userEmail)),
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var user User

	err = DecodeResponseInto(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) AddPendingUser(user *NewUser) (*User, error) {

	body, err := EncodeToJSON(user)
	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/admin/accounts/addpendinguser",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respUser User

	err = DecodeResponseInto(resp, &respUser)
	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

// AddUserToTeamByAdmin - adds user to team with swich account
func (client *Client) AddUserToTeamByAdmin(userID string, accountID string, team string) error {
	// get first accountAdmin and its token
	account, err := client.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	if len(account.Admins) == 0 {
		return fmt.Errorf("Error adding userID %s to Users team of account %s - account does not have any admin", userID, account.Name)
	}

	accountAdminUserID := account.Admins[0]
	accountAdminToken, err := client.GetXAccessToken(accountAdminUserID, accountID)
	if err != nil {
		return err
	}
	// new Client for accountAdmin
	accountAdminClient := NewClient(client.Host, "", accountAdminToken, "x-access-token")
	usersTeam, err := accountAdminClient.GetTeamByName(team)
	if err != nil {
		return err
	}
	if usersTeam == nil {
		fmt.Printf("cannot find users team for account %s", account.Name)
		return nil
	}

	// return is user already assigned to the team
	for _, teamUser := range usersTeam.Users {
		if teamUser.ID == userID {
			return nil
		}
	}

	err = accountAdminClient.AddUserToTeam(usersTeam.ID, userID)

	return err
}

func (client *Client) ActivateUser(userId string) (*User, error) {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/admin/user/%s/activate", userId),
		Method: "POST",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var user User

	err = DecodeResponseInto(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) SetUserAsAccountAdmin(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s/admin", accountId, userId),
		Method: "POST",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteUserAsAccountAdmin(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s/admin", accountId, userId),
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetAllUsers() (*[]User, error) {

	limitPerQuery := 100
	bIsDone := false
	nPageIndex := 1

	var allUsers []User

	for !bIsDone {
		var userPaginatedResp struct {
			Docs []User `json:"docs"`
		}

		opts := RequestOptions{
			Path:   fmt.Sprintf("/admin/user?limit=%d&page=%d", limitPerQuery, nPageIndex),
			Method: "GET",
		}

		resp, err := client.RequestAPI(&opts)

		if err != nil {
			return nil, err
		}

		err = DecodeResponseInto(resp, &userPaginatedResp)

		if err != nil {
			return nil, err
		}

		if len(userPaginatedResp.Docs) > 0 {
			allUsers = append(allUsers, userPaginatedResp.Docs...)
			nPageIndex++
		} else {
			bIsDone = true
		}
	}

	return &allUsers, nil
}

func (client *Client) GetUserByID(userId string) (*User, error) {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/admin/user/id/%s", userId),
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var user User

	err = DecodeResponseInto(resp, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (client *Client) DeleteUser(userName string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/admin/user/%s", userName),
		Method: "DELETE",
	}

	// The API will return a 500 error if the user cannot be found
	// In this case the DeleteUser function should not return an error.
	// Return error only if the body of the return message does not contain "User does not exist"
	res, err := client.RequestAPI(&opts)
	if err != nil {
		if !strings.Contains(string(res), "User does not exist") {
			return err
		}
	}

	return nil
}

func (client *Client) DeleteUserFromAccount(accountId, userId string) error {

	opts := RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/%s", accountId, userId),
		Method: "DELETE",
	}

	_, err := client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func ToSlimAccount(account Account) Account {
	return Account{ID: account.ID}
}

func ToSlimAccounts(accounts []Account) []Account {
	var result []Account
	for i := 0; i < len(accounts); i++ {
		result = append(result, ToSlimAccount(accounts[i]))
	}
	return result
}

func (client *Client) UpdateUserAccounts(userId string, accounts []Account) error {

	// API call '/accounts/{accountId}/{userId}/adduser' doesn't work

	user, err := client.GetUserByID(userId)
	if err != nil {
		return err
	}

	var _accounts = ToSlimAccounts(accounts)

	postUser := UserAccounts{
		UserName: user.UserName,
		Account:  _accounts,
	}

	body, err := EncodeToJSON(postUser)
	if err != nil {
		return err
	}

	opts := RequestOptions{
		Path:   "/admin/user/account",
		Method: "POST",
		Body:   body,
	}

	_, err = client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateUserDetails(accountId, userId, userName, userEmail string) (*User, error) {

	fullPath := fmt.Sprintf("/accounts/%s/%s/updateuser", accountId, userId)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   []byte(generateUserDetailsBody(userName, userEmail)),
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respUser User

	err = DecodeResponseInto(resp, &respUser)
	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (client *Client) UpdateLocalUserPassword(userName, password string) error {

	fullPath := "/admin/user/localProvider"

	requestBody := fmt.Sprintf(`{"userName": "%s","password": "%s"}`, userName, password)

	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   []byte(requestBody),
	}

	_, err := client.RequestAPI(&opts)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteLocalUserPassword(userName string) error {

	fullPath := fmt.Sprintf("/admin/user/localProvider?userName=%s", userName)

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
