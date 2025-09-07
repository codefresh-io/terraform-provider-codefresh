package cfclient

import (
	"fmt"
)

const (
	environmentQueryFields = `
		id
		name
		kind
		clusters {
			runtimeName
			name
			server
			namespaces
		}
		labelPairs
	`
)

// GitopsEnvironment represents a GitOps environment configuration.
type GitopsEnvironment struct {
	ID         string                   `json:"id,omitempty"`
	Name       string                   `json:"name"`
	Kind       string                   `json:"kind"`
	Clusters   []GitopsEnvironmentCluster `json:"clusters"`
	LabelPairs []string                 `json:"labelPairs"`
}

// GitopsCluster represents a cluster within a GitOps environment.
type GitopsEnvironmentCluster struct {
	Name        string   `json:"name"`
	Server      string   `json:"server"`
	RuntimeName string   `json:"runtimeName"`
	Namespaces  []string `json:"namespaces"`
}

type GitopsEnvironmentResponse struct {
	Data struct {
		Environment       GitopsEnvironment `json:"environment,omitempty"`
		CreateEnvironment GitopsEnvironment `json:"createEnvironment,omitempty"`
		UpdateEnvironment GitopsEnvironment `json:"updateEnvironment,omitempty"`
		DeleteEnvironment GitopsEnvironment `json:"deleteEnvironment,omitempty"`
	} `json:"data"`
}

type GitopsEnvironmentsResponse struct {
	Data struct {
		Environments []GitopsEnvironment `json:"environments,omitempty"`
	} `json:"data"`
}

// At the time of writing Codefresh Graphql API does not support fetching environment by ID, So all environemnts are fetched and filtered within this client
func (client *Client) GetGitopsEnvironments() (*[]GitopsEnvironment, error) {

	request := GraphQLRequest{
		Query: fmt.Sprintf(`
			query ($filters: EnvironmentFilterArgs) {
				environments(filters: $filters) {
					%s
				}
			}
		`, environmentQueryFields),
		Variables: map[string]interface{}{},
	}

	response, err := client.SendGqlRequest(request)

	if err != nil {
		return nil, err
	}

	var gitopsEnvironmentsResponse GitopsEnvironmentsResponse
	err = DecodeGraphQLResponseInto(response, &gitopsEnvironmentsResponse)

	if err != nil {
		return nil, err
	}

	return &gitopsEnvironmentsResponse.Data.Environments, nil
}

func (client *Client) GetGitopsEnvironmentById(id string) (*GitopsEnvironment, error) {
	environments, err := client.GetGitopsEnvironments()

	if err != nil {
		return nil, err
	}

	for _, env := range *environments {
		if env.ID == id {
			return &env, nil
		}
	}

	return nil, nil
}

func (client *Client) CreateGitopsEnvironment(environment *GitopsEnvironment) (*GitopsEnvironment, error) {
	request := GraphQLRequest{
		Query: fmt.Sprintf(`
			mutation ($environment: CreateEnvironmentArgs!) {
				createEnvironment(environment: $environment) {
					%s
				}
			}
		`, environmentQueryFields),
		Variables: map[string]interface{}{
			"environment": environment,
		},
	}

	response, err := client.SendGqlRequest(request)

	if err != nil {
		return nil, err
	}

	var gitopsEnvironmentResponse GitopsEnvironmentResponse
	err = DecodeGraphQLResponseInto(response, &gitopsEnvironmentResponse)

	if err != nil {
		return nil, err
	}

	return &gitopsEnvironmentResponse.Data.CreateEnvironment, nil
}

func (client *Client) DeleteGitopsEnvironment(id string) (*GitopsEnvironment, error) {

	type deleteEnvironmentArgs struct {
		ID string `json:"id"`
	}

	request := GraphQLRequest{
		Query: fmt.Sprintf(`
			mutation ($environment: DeleteEnvironmentArgs!) {
  				deleteEnvironment(environment: $environment) {
					%s
				}
			}
		`, environmentQueryFields),

		Variables: map[string]interface{}{
			"environment": deleteEnvironmentArgs{
				ID: id,
			},
		},
	}

	response, err := client.SendGqlRequest(request)

	if err != nil {
		return nil, err
	}
	var gitopsEnvironmentResponse GitopsEnvironmentResponse
	err = DecodeGraphQLResponseInto(response, &gitopsEnvironmentResponse)

	if err != nil {
		return nil, err
	}

	return &gitopsEnvironmentResponse.Data.DeleteEnvironment, nil
}

func (client *Client) UpdateGitopsEnvironment(environment *GitopsEnvironment) (*GitopsEnvironment, error) {
	request := GraphQLRequest{
		Query: fmt.Sprintf(`
			mutation ($environment: UpdateEnvironmentArgs!) {
  				updateEnvironment(environment: $environment) {
					%s
				}
			}
		`, environmentQueryFields),
		Variables: map[string]interface{}{
			"environment": environment,
		},
	}

	response, err := client.SendGqlRequest(request)

	if err != nil {
		return nil, err
	}
	var gitopsEnvironmentResponse GitopsEnvironmentResponse
	err = DecodeGraphQLResponseInto(response, &gitopsEnvironmentResponse)

	if err != nil {
		return nil, err
	}

	return &gitopsEnvironmentResponse.Data.UpdateEnvironment, nil
}
