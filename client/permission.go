package client

import (
	"fmt"
	//"log"
)

// Permission spec
type Permission struct {
	ID       string     `json:"id,omitempty"`
	Team     string     `json:"role,omitempty"`
	Resource string     `json:"resource,omitempty"`
	Action   string     `json:"action,omitempty"`
	Account  string     `json:"account,omitempty"`
	Tags     []string   `json:"attributes,omitempty"`
}

// NewPermission spec, diffs from Permission is `json:"team,omitempty"` vs `json:"role,omitempty"`
type NewPermission struct {
	ID       string     `json:"_id,omitempty"`
	Team     string     `json:"team,omitempty"`
	Resource string     `json:"resource,omitempty"`
	Action   string     `json:"action,omitempty"`
	Account  string     `json:"account,omitempty"`
	Tags     []string   `json:"tags,omitempty"`
}

// GetPermissionList -
func (client *Client) GetPermissionList(teamID, action, resource string) ([]Permission, error) {
	fullPath := "/abac"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var permissions, permissionsFiltered []Permission

	err = DecodeResponseInto(resp, &permissions)
	if err != nil {
		return nil, err
	}

	for _, p := range permissions {
		if teamID != "" && p.Team != teamID {
			continue
		}
		if action != "" && p.Action != action {
			continue
		}
		if resource != "" && p.Resource != resource {
			continue
		}
		permissionsFiltered = append(permissionsFiltered, p)		
	}

	return permissionsFiltered, nil
}

// GetPermissionByID -
func (client *Client) GetPermissionByID(id string) (*Permission, error) {
	fullPath := fmt.Sprintf("/abac/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var permission Permission
	err = DecodeResponseInto(resp, &permission)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

// CreatePermision -
func (client *Client) CreatePermission(permission  *Permission) (*Permission, error) {

	newPermission := &NewPermission{
		ID: permission.ID,
		Team: permission.Team,
		Resource: permission.Resource,
		Action: permission.Action,
		Account: permission.Account,
		Tags: permission.Tags,
	}

	body, err := EncodeToJSON(newPermission)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/abac",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	// respStr := string(resp)
	// log.Printf("[DEBUG] createPermission responce body = %s", respStr)
	var permissionResp []Permission
	err = DecodeResponseInto(resp, &permissionResp)
	if err != nil {
		return nil, err
	}
	if len(permissionResp) != 1 {
		return nil, fmt.Errorf("createPermission - unknown response lenght!=1:  %v", permissionResp)
	}

	newPermissionID := permissionResp[0].ID
    
	return client.GetPermissionByID(newPermissionID)
}

// DeletePermission -
func (client *Client) DeletePermission(id string) error {
	fullPath := fmt.Sprintf("/abac/%s", id)
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