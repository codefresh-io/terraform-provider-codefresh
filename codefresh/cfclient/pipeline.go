package cfclient

import (
	"errors"
	"fmt"
	"strings"
)

type Pipelines struct {
	Docs  []Pipeline `json:"docs,omitempty"`
	Count int        `json:"count,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Labels struct {
	Tags []string `json:"tags,omitempty"`
}

type Metadata struct {
	Name               string `json:"name,omitempty"`
	ID                 string `json:"id,omitempty"`
	IsPublic           bool   `json:"isPublic,omitempty"`
	Labels             Labels `json:"labels,omitempty"`
	OriginalYamlString string `json:"originalYamlString,omitempty"`
	Project            string `json:"project,omitempty"`
	ProjectId          string `json:"projectId,omitempty"`
	Revision           int    `json:"revision,omitempty"`
}

type SpecTemplate struct {
	Location string `json:"location,omitempty"`
	Repo     string `json:"repo,omitempty"`
	Path     string `json:"path,omitempty"`
	Revision string `json:"revision,omitempty"`
	Context  string `json:"context,omitempty"`
}

type Trigger struct {
	Name                         string              `json:"name,omitempty"`
	Description                  string              `json:"description,omitempty"`
	Type                         string              `json:"type,omitempty"`
	Repo                         string              `json:"repo,omitempty"`
	Events                       []string            `json:"events,omitempty"`
	BranchRegex                  string              `json:"branchRegex,omitempty"`
	BranchRegexInput             string              `json:"branchRegexInput,omitempty"`
	PullRequestTargetBranchRegex string              `json:"pullRequestTargetBranchRegex,omitempty"`
	CommentRegex                 string              `json:"commentRegex,omitempty"`
	ModifiedFilesGlob            string              `json:"modifiedFilesGlob,omitempty"`
	Provider                     string              `json:"provider,omitempty"`
	Disabled                     bool                `json:"disabled,omitempty"`
	Options                      *TriggerOptions     `json:"options,omitempty"`
	PullRequestAllowForkEvents   bool                `json:"pullRequestAllowForkEvents,omitempty"`
	CommitStatusTitle            string              `json:"commitStatusTitle,omitempty"`
	Context                      string              `json:"context,omitempty"`
	Contexts                     []string            `json:"contexts,omitempty"`
	RuntimeEnvironment           *RuntimeEnvironment `json:"runtimeEnvironment,omitempty"`
	Variables                    []Variable          `json:"variables,omitempty"`
}

type CronTrigger struct {
	Name               string              `json:"name,omitempty"`
	Type               string              `json:"type,omitempty"`
	Expression         string              `json:"expression,omitempty"`
	Message            string              `json:"message,omitempty"`
	GitTriggerId       string              `json:"gitTriggerId,omitempty"`
	Branch             string              `json:"branch,omitempty"`
	Disabled           bool                `json:"disabled,omitempty"`
	Options            *TriggerOptions     `json:"options,omitempty"`
	RuntimeEnvironment *RuntimeEnvironment `json:"runtimeEnvironment,omitempty"`
	Variables          []Variable          `json:"variables,omitempty"`
}

type TriggerOptions struct {
	NoCache             bool `json:"noCache,omitempty"`
	NoCfCache           bool `json:"noCfCache,omitempty"`
	ResetVolume         bool `json:"resetVolume,omitempty"`
	EnableNotifications bool `json:"enableNotifications,omitempty"`
}

type RuntimeEnvironment struct {
	Name                     string `json:"name,omitempty"`
	Memory                   string `json:"memory,omitempty"`
	CPU                      string `json:"cpu,omitempty"`
	DindStorage              string `json:"dindStorage,omitempty"`
	RequiredAvailableStorage string `json:"requiredAvailableStorage,omitempty"`
}

type ExternalResource struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	Context     string `json:"context"`
	Destination string `json:"destination"`
	IsFolder    bool   `json:"isFolder"`
	Repo        string `json:"repo"`
	Revision    string `json:"revision"`
}

func (t *Trigger) SetVariables(variables map[string]interface{}, encrypted bool) {
	for key, value := range variables {
		t.Variables = append(t.Variables, Variable{Key: key, Value: value.(string), Encrypted: encrypted})
	}
}

func (t *CronTrigger) SetVariables(variables map[string]interface{}, encrypted bool) {
	for key, value := range variables {
		t.Variables = append(t.Variables, Variable{Key: key, Value: value.(string), Encrypted: encrypted})
	}
}

type Spec struct {
	Variables                    []Variable               `json:"variables,omitempty"`
	SpecTemplate                 *SpecTemplate            `json:"specTemplate,omitempty"`
	Triggers                     []Trigger                `json:"triggers,omitempty"`
	CronTriggers                 []CronTrigger            `json:"cronTriggers,omitempty"`
	Priority                     int                      `json:"priority,omitempty"`
	Concurrency                  int                      `json:"concurrency,omitempty"`
	BranchConcurrency            int                      `json:"branchConcurrency,omitempty"`
	TriggerConcurrency           int                      `json:"triggerConcurrency,omitempty"`
	Contexts                     []interface{}            `json:"contexts,omitempty"`
	Steps                        *Steps                   `json:"steps,omitempty"`
	Stages                       *Stages                  `json:"stages,omitempty"`
	Mode                         string                   `json:"mode,omitempty"`
	FailFast                     *bool                    `json:"fail_fast,omitempty"`
	RuntimeEnvironment           RuntimeEnvironment       `json:"runtimeEnvironment,omitempty"`
	TerminationPolicy            []map[string]interface{} `json:"terminationPolicy,omitempty"`
	PackId                       string                   `json:"packId,omitempty"`
	RequiredAvailableStorage     string                   `json:"requiredAvailableStorage,omitempty"`
	Hooks                        *Hooks                   `json:"hooks,omitempty"`
	Options                      map[string]bool          `json:"options,omitempty"`
	PermitRestartFromFailedSteps *bool                    `json:"permitRestartFromFailedSteps,omitempty"`
	ExternalResources            []ExternalResource       `json:"externalResources,omitempty"`
}

type Steps struct {
	Steps string
}
type Stages struct {
	Stages string
}

type Hooks struct {
	Hooks string
}

func (d Steps) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Steps)
	return bytes, nil
}
func (d Stages) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Stages)
	return bytes, nil
}

func (d Hooks) MarshalJSON() ([]byte, error) {
	bytes := []byte(d.Hooks)
	return bytes, nil
}
func (d *Steps) UnmarshalJSON(data []byte) error {
	d.Steps = string(data)
	return nil
}
func (d *Stages) UnmarshalJSON(data []byte) error {
	d.Stages = string(data)
	return nil
}
func (d *Hooks) UnmarshalJSON(data []byte) error {
	d.Hooks = string(data)
	return nil
}

type Pipeline struct {
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     Spec     `json:"spec,omitempty"`
	Version  string   `json:"version,omitempty"`
}

func (p *Pipeline) SetVariables(variables map[string]interface{}, encrypted bool) {
	for key, value := range variables {
		p.Spec.Variables = append(p.Spec.Variables, Variable{Key: key, Value: value.(string), Encrypted: encrypted})
	}
}

func (pipeline *Pipeline) GetID() string {
	if pipeline.Metadata.ID != "" {
		return pipeline.Metadata.ID
	} else {
		return pipeline.Metadata.Name
	}
}

func (client *Client) GetPipeline(name string) (*Pipeline, error) {
	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(name, "/", "%2F", 1))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var pipeline Pipeline

	err = DecodeResponseInto(resp, &pipeline)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (client *Client) GetPipelines() (*[]Pipeline, error) {
	fullPath := "/pipelines"
	opts := RequestOptions{
		Path:   fullPath,
		Method: "GET",
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var getPipelines Pipelines

	err = DecodeResponseInto(resp, &getPipelines)
	if err != nil {
		return nil, err
	}

	return &getPipelines.Docs, nil
}

func (client *Client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {

	body, err := EncodeToJSON(pipeline)

	if err != nil {
		return nil, err
	}
	opts := RequestOptions{
		Path:   "/pipelines",
		Method: "POST",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)

	if err != nil {
		return nil, err
	}

	var respPipeline Pipeline
	err = DecodeResponseInto(resp, &respPipeline)
	if err != nil {
		return nil, err
	}

	return &respPipeline, nil

}

func (client *Client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {

	body, err := EncodeToJSON(pipeline)

	if err != nil {
		return nil, err
	}

	id := pipeline.GetID()
	if id == "" {
		return nil, errors.New("[ERROR] Both Pipeline ID and Name are empty")
	}

	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(id, "/", "%2F", 1))
	opts := RequestOptions{
		Path:   fullPath,
		Method: "PUT",
		Body:   body,
	}

	resp, err := client.RequestAPI(&opts)
	if err != nil {
		return nil, err
	}

	var respPipeline Pipeline
	err = DecodeResponseInto(resp, &respPipeline)
	if err != nil {
		return nil, err
	}

	return &respPipeline, nil
}

func (client *Client) DeletePipeline(name string) error {

	fullPath := fmt.Sprintf("/pipelines/%s", strings.Replace(name, "/", "%2F", 1))
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
