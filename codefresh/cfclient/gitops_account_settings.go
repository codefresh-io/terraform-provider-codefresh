package cfclient

import (
	"fmt"
)

type GitopsActiveAccountResponse struct {
	Data struct {
		Me struct {
			ActiveAccount GitopsActiveAccountInfo `json:"activeAccount,omitempty"`
		} `json:"me,omitempty"`
	} `json:"data,omitempty"`
}

type GitopsActiveAccountInfo struct {
	ID               string   `json:"id,omitempty"`
	AccountName      string   `json:"name,omitempty"`
	GitProvider      string   `json:"gitProvider,omitempty"`
	GitApiUrl        string   `json:"gitApiUrl,omitempty"`
	SharedConfigRepo string   `json:"sharedConfigRepo,omitempty"`
	Admins           []string `json:"admins,omitempty"`
}

func (client *Client) GetActiveGitopsAccountInfo() (*GitopsActiveAccountInfo, error) {
	request := GraphQLRequest{
		Query: `
			query AccountInfo {
				me {
					activeAccount {
						id
						name
						gitProvider
						gitApiUrl
						sharedConfigRepo
						admins
					}
				}
			}
		`,
	}

	response, err := client.SendGqlRequest(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var gitopsAccountResponse GitopsActiveAccountResponse

	err = DecodeGraphQLResponseInto(response, &gitopsAccountResponse)

	if err != nil {
		return nil, err
	}

	gitopsActiveAccountInfo := gitopsAccountResponse.Data.Me.ActiveAccount

	return &gitopsActiveAccountInfo, nil
}
