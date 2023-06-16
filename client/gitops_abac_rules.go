package client

import (
	"fmt"
)

type AbacEntityValues string

const (
	AbacEntityClusters            AbacEntityValues = "clusters"
	AbacEntityExecutionContext    AbacEntityValues = "executionContext"
	AbacEntityGitContexts         AbacEntityValues = "gitContexts"
	AbacEntityGitopsApplications  AbacEntityValues = "gitopsApplications"
	AbacEntityHelmCharts          AbacEntityValues = "helmCharts"
	AbacEntityPipelines           AbacEntityValues = "pipelines"
	AbacEntityProjects            AbacEntityValues = "projects"
	AbacEntitySharedConfiguration AbacEntityValues = "sharedConfiguration"
	AbacEntityWorkflows           AbacEntityValues = "workflows"
	AbacEntityWorkflowTemplates   AbacEntityValues = "workflowTemplates"
)

type EntityAbacAttribute struct {
	Name  string `json:"name"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value"`
}

// GitopsAbacRule spec
type GitopsAbacRule struct {
	ID         string                `json:"id,omitempty"`
	AccountId  string                `json:"accountId,omitempty"`
	EntityType AbacEntityValues      `json:"entityType"`
	Teams      []string              `json:"teams"`
	Tags       []string              `json:"tags,omitempty"`
	Actions    []string              `json:"actions"`
	Attributes []EntityAbacAttribute `json:"attributes"`
}

type GitopsAbacRulesListResponse struct {
	Data struct {
		AbacRules []GitopsAbacRule `json:"abacRules"`
	} `json:"data"`
}

type GitopsAbacRuleResponse struct {
	Data struct {
		AbacRule       GitopsAbacRule `json:"abacRule,omitempty"`
		CreateAbacRule GitopsAbacRule `json:"createAbacRule,omitempty"`
		RemoveAbacRule GitopsAbacRule `json:"removeAbacRule,omitempty"`
	} `json:"data"`
}

func (client *Client) GetAbacRulesList(accountId string, entityType AbacEntityValues) ([]GitopsAbacRule, error) {
	request := GraphQLRequest{
		Query: `
			query AbacRules($accountId: String!, $entityType: AbacEntityValues!) {
				abacRules(accountId: $accountId, entityType: $entityType) {
					id
					accountId
					entityType
					teams
					tags
					actions
					attributes {
						name
						key
						value
					}
				}
			}
		`,
		Variables: map[string]interface{}{
			"accountId":  accountId,
			"entityType": entityType,
		},
	}

	response, err := client.SendGqlRequest(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var gitopsAbacRulesResponse GitopsAbacRulesListResponse
	err = DecodeGraphQLResponseInto(response, &gitopsAbacRulesResponse)
	if err != nil {
		return nil, err
	}

	return gitopsAbacRulesResponse.Data.AbacRules, nil
}

// GetAbacRuleByID -
func (client *Client) GetAbacRuleByID(accountId string, id string) (*GitopsAbacRule, error) {
	request := GraphQLRequest{
		Query: `
			query AbacRule($accountId: String!, $id: ID!) {
				abacRule(accountId: $accountId, id: $id) {
					id
					accountId
					entityType
					teams
					tags
					actions
					attributes {
						name
						key
						value
					}
				}
			}
		`,
		Variables: map[string]interface{}{
			"accountId": accountId,
			"id":        id,
		},
	}

	response, err := client.SendGqlRequest(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var gitopsAbacRuleResponse GitopsAbacRuleResponse
	err = DecodeGraphQLResponseInto(response, &gitopsAbacRuleResponse)
	if err != nil {
		return nil, err
	}

	return &gitopsAbacRuleResponse.Data.AbacRule, nil
}

func (client *Client) CreateAbacRule(accountId string, gitopsAbacRule *GitopsAbacRule) (*GitopsAbacRule, error) {

	newAbacRule := &GitopsAbacRule{
		EntityType: gitopsAbacRule.EntityType,
		Teams:      gitopsAbacRule.Teams,
		Tags:       gitopsAbacRule.Tags,
		Actions:    gitopsAbacRule.Actions,
		Attributes: gitopsAbacRule.Attributes,
	}

	request := GraphQLRequest{
		Query: `
			mutation CreateAbacRule($accountId: String!, $createAbacRuleInput: CreateAbacRuleInput!) {
				createAbacRule(accountId: $accountId, createAbacRuleInput: $createAbacRuleInput) {
					id
					accountId
					entityType
					teams
					tags
					actions
					attributes {
						name
						key
						value
					}
				}
			}
		`,
		Variables: map[string]interface{}{
			"accountId":           accountId,
			"createAbacRuleInput": newAbacRule,
		},
	}

	response, err := client.SendGqlRequest(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var gitopsAbacRuleResponse GitopsAbacRuleResponse
	err = DecodeGraphQLResponseInto(response, &gitopsAbacRuleResponse)
	if err != nil {
		return nil, err
	}

	return &gitopsAbacRuleResponse.Data.CreateAbacRule, nil
}

func (client *Client) DeleteAbacRule(accountId string, id string) (*GitopsAbacRule, error) {
	request := GraphQLRequest{
		Query: `
			mutation RemoveAbacRule($accountId: String!, $id: ID!) {
				removeAbacRule(accountId: $accountId, id: $id) {
					id
					accountId
					entityType
					teams
					tags
					actions
					attributes {
						name
						key
						value
					}
				}
			}
		`,
		Variables: map[string]interface{}{
			"accountId": accountId,
			"id":        id,
		},
	}

	response, err := client.SendGqlRequest(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var gitopsAbacRuleResponse GitopsAbacRuleResponse
	err = DecodeGraphQLResponseInto(response, &gitopsAbacRuleResponse)
	if err != nil {
		return nil, err
	}

	return &gitopsAbacRuleResponse.Data.RemoveAbacRule, nil
}
