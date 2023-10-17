package cfclient

import (
	"errors"
	"fmt"
)

// Project spec
type Project struct {
	ID          string     `json:"id,omitempty"`
	ProjectName string     `json:"projectName,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Variables   []Variable `json:"variables,omitempty"`
}

// GetID implement CodefreshObject interface
func (project *Project) GetID() string {
	return project.ID
}

// SetVariables project variables
func (project *Project) SetVariables(variables map[string]interface{}) {
	for key, value := range variables {
		project.Variables = append(project.Variables, Variable{Key: key, Value: value.(string)})
	}
}

// GetProjectByName get project object by name
func (client *Client) GetProjectByName(name string) (*Project, error) {
	fullPath := fmt.Sprintf("/projects/name/%s", name)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var project Project

	err = DecodeResponseInto(resp, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProjectByID get project object by id
func (client *Client) GetProjectByID(id string) (*Project, error) {
	fullPath := fmt.Sprintf("/projects/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var project Project

	err = DecodeResponseInto(resp, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject POST project
func (client *Client) CreateProject(project *Project) (*Project, error) {

	body, err := EncodeToJSON(project)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/projects",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respProject Project
	err = DecodeResponseInto(resp, &respProject)
	if err != nil {
		return nil, err
	}

	return &respProject, nil
}

// UpdateProject PATCH project
func (client *Client) UpdateProject(project *Project) error {

	body, err := EncodeToJSON(project)

	if err != nil {
		return err
	}

	id := project.GetID()
	if id == "" {
		return errors.New("[ERROR] Project ID is empty")
	}

	fullPath := fmt.Sprintf("/projects/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PATCH",
		Body:   body,
	}

	_, err = client.RequestAPI(&opts)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProject DELETE
func (client *Client) DeleteProject(id string) error {
	fullPath := fmt.Sprintf("/projects/%s", id)
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
