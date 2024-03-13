package cfclient

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stretchr/objx"
)

// CurrentAccountUser spec
type CurrentAccountUser struct {
	ID       string `json:"id,omitempty"`
	UserName string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Status   string `json:"status,omitempty"`
}

// CurrentAccount spec
type CurrentAccount struct {
	ID     string
	Name   string
	Users  []CurrentAccountUser
	Admins []CurrentAccountUser
}

// GetCurrentAccount -
func (client *Client) GetCurrentAccount() (*CurrentAccount, error) {

	// get and parse current account
	userResp, err := client.RequestAPI(&RequestOptions{
		Path:   "/user",
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}
	userRespStr := string(userResp)
	currentAccountX, err := objx.FromJSON(userRespStr)
	if err != nil {
		return nil, err
	}

	activeAccountName := currentAccountX.Get("activeAccountName").String()
	if activeAccountName == "" {
		return nil, fmt.Errorf("GetCurrentAccount - cannot get activeAccountName")
	}
	currentAccount := &CurrentAccount{
		Name:   activeAccountName,
		Users:  make([]CurrentAccountUser, 0),
		Admins: make([]CurrentAccountUser, 0),
	}

	accountAdminsIDs := make([]string, 0)

	allAccountsI := currentAccountX.Get("account").InterSlice()
	for _, accI := range allAccountsI {
		accX := objx.New(accI)
		if accX.Get("name").String() == activeAccountName {
			currentAccount.ID = accX.Get("id").String()
			admins := accX.Get("admins").InterSlice()
			for _, adminI := range admins {
				accountAdminsIDs = append(accountAdminsIDs ,adminI.(string))
			}
			break
		}
	}
	if currentAccount.ID == "" {
		return nil, fmt.Errorf("GetCurrentAccount - cannot get activeAccountName")
	}

	// get and parse account users
	accountUsersResp, err := client.RequestAPI(&RequestOptions{
		Path:   fmt.Sprintf("/accounts/%s/users", currentAccount.ID),
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	accountUsersI := make([]interface{}, 0)
	if e := json.Unmarshal(accountUsersResp, &accountUsersI); e != nil {
		return nil, fmt.Errorf("cannot unmarshal accountUsers responce for accountId=%s: %v", currentAccount.ID, e)
	}
	for _, userI := range accountUsersI {
		userX := objx.New(userI)
		userName := userX.Get("userName").String()
		email := userX.Get("email").String()
		status := userX.Get("status").String()
		userID := userX.Get("_id").String()
		currentAccount.Users = append(currentAccount.Users, CurrentAccountUser{
			ID:       userID,
			UserName: userName,
			Email:    email,
			Status:   status,
		})

		// If user exists in Admin list append it to addmins as well. This assumes that a user cannot be an admin without being a regular user too.
		if slices.Contains(accountAdminsIDs, userID) {
			currentAccount.Admins = append(currentAccount.Admins, CurrentAccountUser{
				ID:       userID,
				UserName: userName,
				Email:    email,
				Status:   status,
			})
		}
	}

	return currentAccount, nil
}
