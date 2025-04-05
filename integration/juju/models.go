package juju

import (
	"encoding/json"
	"fmt"
)

const (
	addModelCommand   = "add-model"
	listModelsCommand = "list-models"
)

type AddModelOptions struct {
	Name string
}

type Credential struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Cloud  string `json:"cloud"`
	Type   string `json:"type"`
	Region string `json:"region"`
}

type ModelStatus struct {
	Name    string `json:"status"`
	Message string `json:"message"`
	Current string `json:"current"`
	Since   string `json:"since"`
}

type Model struct {
	Name           string      `json:"name"`
	ShortName      string      `json:"short-name"`
	ModelUUID      string      `json:"model-uuid"`
	ModelType      string      `json:"model-type"`
	ControllerUUID string      `json:"controller-uuid"`
	ControllerName string      `json:"controller-name"`
	IsController   bool        `json:"is-controller"`
	Owner          string      `json:"owner"`
	Cloud          string      `json:"cloud"`
	Region         string      `json:"region"`
	Credential     Credential  `json:"credential"`
	Type           string      `json:"type"`
	Life           string      `json:"life"`
	ModelStatus    ModelStatus `json:"status"`
	Access         string      `json:"access"`
	LastConnection string      `json:"last-connection"`
	SlaOwner       string      `json:"sla-owner"`
	AgentVersion   string      `json:"agent-version"`
}

type ListModelOutput struct {
	Models []*Model `json:"models"`
}

func (j Client) AddModel(opts *AddModelOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("model name is required")
	}

	args := []string{addModelCommand, opts.Name}

	_, err := j.Runner.Run(args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

func (j Client) ListModels() ([]*Model, error) {
	args := []string{listModelsCommand, "--format=json"}

	output, err := j.Runner.Run(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	var listModelOutput *ListModelOutput

	err = json.Unmarshal(output, &listModelOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal models: %w", err)
	}

	if listModelOutput == nil {
		return nil, nil
	}

	return listModelOutput.Models, nil
}
