package client

import "fmt"

type User struct {
	ID             string              `json:"_id"`
	UserName       string              `json:"userName"`
	Email          string              `json:"email"`
	Roles          []interface{}       `json:"roles"`
	DefaultAccount int                 `json:"defaultAccount"`
	Account        []Account           `json:"account"`
	Status         string              `json:"status"`
	RegisterDate   string              `json:"register_date"`
	HasPassword    bool                `json:"hasPassword"`
	Notifications  []NotificationEvent `json:"notifications"`
	ShortProfile   struct {
		UserName string `json:"userName"`
	} `json:"shortProfile"`
	Settings struct {
		SendWeeklyReport bool `json:"sendWeeklyReport"`
	} `json:"settings"`
	Logins    []interface{} `json:"logins"`
	InviteURL string        `json:"inviteUrl"`
}

func (client *Client) AddNewUserToAccount(accountId, userName, userEmail string) (*User, error) {

	userDetails := fmt.Sprintf(`{"userName": "%s", "email": "%s"}`, userName, userEmail)

	fullPath := fmt.Sprintf("/accounts/%s/adduser", accountId)

	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   []byte(userDetails),
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