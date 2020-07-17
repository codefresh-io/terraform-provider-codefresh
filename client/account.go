package client

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/imdario/mergo"
)

type DockerRegistry struct {
	Kind                string `json:"kind"`
	BehindFirewall      bool   `json:"behindFirewall"`
	Primary             bool   `json:"primary"`
	Default             bool   `json:"default"`
	Internal            bool   `json:"internal"`
	DenyCompositeDomain bool   `json:"denyCompositeDomain"`
	ID                  string `json:"_id"`
	Name                string `json:"name"`
	Provider            string `json:"provider"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Domain              string `json:"domain"`
	RepositoryPrefix    string `json:"repositoryPrefix"`
}

type Suspension struct {
	IsSuspended bool `json:"isSuspended"`
}

type Exist struct {
	Exist bool `json:"exist"`
}

type Active struct {
	Active bool `json:"active"`
}

type Integration struct {
	Stash        Active           `json:"stash,omitempty"`
	Github       Active           `json:"github,omitempty"`
	Gitlab       Active           `json:"gitlab,omitempty"`
	Aks          Exist            `json:"aks,omitempty"`
	Gcloud       Exist            `json:"gcloud,omitempty"`
	DigitalOcean Exist            `json:"digitalOcean,omitempty"`
	Registries   []DockerRegistry `json:"registries,omitempty"`
}

type Build struct {
	Strategy string        `json:"strategy,omitempty"`
	Nodes    int           `json:"nodes,omitempty"`
	Parallel int           `json:"parallel,omitempty"`
	Packs    []interface{} `json:"packs,omitempty"`
}

type PaymentPlan struct {
	Trial struct {
		Trialing             bool   `json:"trialing"`
		TrialWillEndNotified bool   `json:"trialWillEndNotified"`
		TrialEndedNotified   bool   `json:"trialEndedNotified"`
		Type                 string `json:"type"`
		PreviousSegment      string `json:"previousSegment"`
	} `json:"trial,omitempty"`
	ID       string `json:"id,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type ImageViewConfig struct {
	Version string `json:"version"`
}

type BuildStepConfig struct {
	Version     string `json:"version"`
	DisablePush bool   `json:"disablePush"`
	AutoPush    bool   `json:"autoPush"`
}

type CFCRState struct {
	Enabled             bool   `json:"enabled"`
	System              string `json:"system"`
	DisplayGlobalNotice bool   `json:"displayGlobalNotice"`
	AccountChoice       string `json:"accountChoice"`
}

type NotificationEvent struct {
	Events []string `json:"events"`
	Type   string   `json:"type"`
}

type Collaborators struct {
	Limit int `json:"limit"`
	Used  int `json:"used,omitempty"`
}

type DataRetention struct {
	Weeks int `json:"weeks"`
}

type Limits struct {
	Collaborators Collaborators `json:"collaborators,omitempty"`
	DataRetention DataRetention `json:"dataRetention,omitempty"`
}

type Account struct {
	Suspension                  *Suspension         `json:"suspension,omitempty"`
	Integrations                *Integration        `json:"integrations,omitempty"`
	Build                       *Build              `json:"build,omitempty"`
	PaymentPlan                 *PaymentPlan        `json:"paymentPlan,omitempty"`
	ImageViewConfig             *ImageViewConfig    `json:"imageViewConfig,omitempty"`
	BuildStepConfig             *BuildStepConfig    `json:"buildStepConfig,omitempty"`
	CFCRState                   *CFCRState          `json:"CFCRState,omitempty"`
	AllowedDomains              []interface{}       `json:"allowedDomains,omitempty"`
	Admins                      []string            `json:"admins,omitempty"`
	Environment                 int                 `json:"environment,omitempty"`
	DedicatedInfrastructure     bool                `json:"dedicatedInfrastructure,omitempty"`
	CanUsePrivateRepos          bool                `json:"canUsePrivateRepos,omitempty"`
	SupportPlan                 string              `json:"supportPlan,omitempty"`
	IncreasedAttention          bool                `json:"increasedAttention,omitempty"`
	LocalUserPasswordIDPEnabled bool                `json:"localUserPasswordIDPEnabled,omitempty"`
	CodefreshEnv                string              `json:"codefreshEnv,omitempty"`
	ID                          string              `json:"_id,omitempty"`
	BadgeToken                  string              `json:"badgeToken,omitempty"`
	CreatedAt                   string              `json:"createdAt,omitempty"`
	UpdatedAt                   string              `json:"updatedAt,omitempty"`
	Name                        string              `json:"name,omitempty"`
	RuntimeEnvironment          string              `json:"runtimeEnvironment,omitempty"`
	CfcrRepositoryPath          string              `json:"cfcrRepositoryPath,omitempty"`
	Notifications               []NotificationEvent `json:"notifications,omitempty"`
	RepoPermission              string              `json:"repoPermission,omitempty"`
	Limits                      *Limits             `json:"limits,omitempty"`
	Features                    map[string]bool     `json:"features,omitempty"`
	// Features                    *Features           `json:"features,omitempty"`
	// RuntimeEnvironments ToDo
	// Remaining ToDo
	// ID string `json:"id"`
}

type AccountDetails struct {
	AccountDetails Account `json:"accountDetails"`
}

// Decodes a TypeMap of map[string]interface{} into map[string]bool for account features
func (account *Account) SetFeatures(m map[string]interface{}) {
	res := make(map[string]bool)
	for k, v := range m {
		value := v.(string)
		b, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("[ERROR] Can't parse '%s = %s' as boolean", k, value)
		}
		res[k] = b
	}
	account.Features = res
}

func (account *Account) GetID() string {
	return account.ID
}

func (client *Client) GetAccountByID(id string) (*Account, error) {
	fullPath := fmt.Sprintf("/admin/accounts/%s", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var account Account

	err = DecodeResponseInto(resp, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (client *Client) GetAllAccounts() (*[]Account, error) {

	opts := RequestOptions{
		Path:   "/admin/accounts",
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var accounts []Account

	err = DecodeResponseInto(resp, &accounts)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (client *Client) GetAccountsList(accountsId []string) (*[]Account, error) {

	var accounts []Account

	for _, accountId := range accountsId {
		account, err := client.GetAccountByID(accountId)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, *account)
	}

	return &accounts, nil
}

func (client *Client) CreateAccount(account *Account) (*Account, error) {

	body, err := EncodeToJSON(account)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/admin/accounts",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respAccount Account

	err = DecodeResponseInto(resp, &respAccount)
	if err != nil {
		return nil, err
	}

	return &respAccount, nil
}

func (client *Client) UpdateAccount(account *Account) (*Account, error) {

	id := account.GetID()
	if id == "" {
		return nil, errors.New("[ERROR] Account ID is empty")
	}

	existingAccount, err := client.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(account, existingAccount)
	if err != nil {
		return nil, err
	}

	putAccount := AccountDetails{*account}

	body, err := EncodeToJSON(putAccount)
	if err != nil {
		return nil, err
	}

	fullPath := fmt.Sprintf("/admin/accounts/%s/update", id)
	opts := RequestOptions{
		Path:   fullPath,
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respAccount Account
	err = DecodeResponseInto(resp, &respAccount)
	if err != nil {
		return nil, err
	}

	return &respAccount, nil
}

func (client *Client) DeleteAccount(id string) error {

	fullPath := fmt.Sprintf("/admin/accounts/%s", id)
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

func GetAccountAdminsDiff(desiredAdmins []string, existingAdmins []string) (adminsToAdd []string, adminsToDelete []string) {

	adminsToAdd = []string{}
	adminsToDelete = []string{}

	for _, id := range existingAdmins {
		if ok := FindInSlice(desiredAdmins, id); !ok {
			adminsToDelete = append(adminsToDelete, id)
		}
	}

	for _, id := range desiredAdmins {

		if ok := FindInSlice(existingAdmins, id); !ok {
			adminsToAdd = append(adminsToAdd, id)
		}
	}

	return adminsToAdd, adminsToDelete
}
