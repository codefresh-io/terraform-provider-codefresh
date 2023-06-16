package client_test

import (
	"fmt"
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/codefresh-io/terraform-provider-codefresh/codefresh"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Init tests

	// Run tests
	exitCode := m.Run()

	// Cleanup

	// Exit
	os.Exit(exitCode)
}

// TODO: Create and remove account, team
func TestRules(t *testing.T) {
	// Validate environment variables
	url := os.Getenv("CODEFRESH_API_URL")
	if url == "" {
		t.Fatalf("CODEFRESH_API_URL variable is not set")
	}
	url2 := os.Getenv("CODEFRESH_API2_URL")
	if url2 == "" {
		t.Fatalf("CODEFRESH_API2_URL variable is not set")
	}
	key := os.Getenv("CODEFRESH_API_KEY")
	if key == "" {
		t.Fatalf("CODEFRESH_API_KEY variable is not set")
	}

	client := cfClient.NewClient(
		os.Getenv(codefresh.ENV_CODEFRESH_API_URL),
		os.Getenv(codefresh.ENV_CODEFRESH_API2_URL),
		os.Getenv(codefresh.ENV_CODEFRESH_API_KEY),
		"",
	)

	currentAccount, err := client.GetCurrentAccount()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(currentAccount)

	created, err := client.CreateAbacRule(
		currentAccount.ID,
		&cfClient.GitopsAbacRule{
			EntityType: cfClient.AbacEntityGitopsApplications,
			Teams:      []string{"6365495094c782ba1ba45451"},
			Tags:       []string{},
			Actions:    []string{"SYNC"},
			Attributes: []cfClient.EntityAbacAttribute{
				{
					Name:  "LABEL",
					Key:   "AnyName",
					Value: "SomeValue",
				},
				{
					Name:  "NAMESPACE",
					Value: "local51",
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(created)
	if created == nil {
		t.Fatalf("Empty rule after creation")
	}

	list, err := client.GetAbacRulesList(currentAccount.ID, cfClient.AbacEntityGitopsApplications)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(list)
	if len(list) == 0 {
		t.Fatalf("List of rules is empty")
	}

	one, err := client.GetAbacRuleByID(currentAccount.ID, created.ID)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(one)
	if created.ID != one.ID {
		t.Fatalf("Expected %s, but got %s", created.ID, one.ID)
	}

	deleted, err := client.DeleteAbacRule(currentAccount.ID, created.ID)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}
	fmt.Println(deleted)
	if created.ID != deleted.ID {
		t.Fatalf("Expected %s, but got %s", created.ID, deleted.ID)
	}
}
