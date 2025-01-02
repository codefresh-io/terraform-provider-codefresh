package cfclient

import (
	"fmt"
	"golang.org/x/exp/slices"

)

type ServiceUser struct {
	ID       	   string   `json:"_id,omitempty"`
	Name 	 	   string   `json:"userName,omitempty"`
	Teams  	       []Team   `json:"teams,omitempty"`
	Roles 		   []string `json:"roles,omitempty"`
}

type ServiceUserCreateUpdate struct {
	ID       	    string   `json:"_id,omitempty"`
	Name 	 	    string   `json:"userName,omitempty"`
	TeamIDs		    []string `json:"teamIds,omitempty"`
	AssignAdminRole bool    `json:"assignAdminRole,omitempty"`
}

// GetID implement CodefreshObject interface
func (serviceuser *ServiceUser) GetID() string {
	return serviceuser.ID
}

func (serviceuser *ServiceUser) HasAdminRole() bool {
	return slices.Contains(serviceuser.Roles, "Admin")
}

func (client *Client) GetServiceUserList() ([]ServiceUser, error) {
	fullPath := "/service-users"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var serviceusers []ServiceUser

	err = DecodeResponseInto(resp, &serviceusers)
	if err != nil {
		return nil, err
	}

	return serviceusers, nil
}

func (client *Client) GetServiceUserByName(name string) (*ServiceUser, error) {

	serviceusers, err := client.GetServiceUserList()
	if err != nil {
		return nil, err
	}

	for _, serviceuser := range serviceusers {
		if serviceuser.Name == name {
			return &serviceuser, nil
		}
	}

	return nil, nil
}

func (client *Client) GetServiceUserByID(id string) (*ServiceUser, error) {

	fullPath := fmt.Sprintf("/service-users/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var serviceuser ServiceUser

	err = DecodeResponseInto(resp, &serviceuser)
	if err != nil {
		return nil, err
	}

	return &serviceuser, nil
}

func (client *Client) CreateServiceUser(serviceUserCreateUpdate *ServiceUserCreateUpdate) (*ServiceUser, error) {

	fullPath := "/service-users"
	body, err := EncodeToJSON(serviceUserCreateUpdate)

	if err != nil {
		return nil, err
	}

	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body: body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var serviceuser ServiceUser

	err = DecodeResponseInto(resp, &serviceuser)
	if err != nil {
		return nil, err
	}

	return &serviceuser, nil
}

func (client *Client) UpdateServiceUser(serviceUserCreateUpdate *ServiceUserCreateUpdate) (*ServiceUser, error) {

	fullPath := fmt.Sprintf("/service-users/%s", serviceUserCreateUpdate.ID)
	body, err := EncodeToJSON(serviceUserCreateUpdate)

	if err != nil {
		return nil, err
	}

	opts := RequestOptions{
		Path:   fullPath,
		Method: "PATCH",
		Body: body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var serviceuser ServiceUser

	err = DecodeResponseInto(resp, &serviceuser)
	if err != nil {
		return nil, err
	}

	return &serviceuser, nil
}

func (client *Client) DeleteServiceUser(id string) error {
	fullPath := fmt.Sprintf("/service-users/%s", id)
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
