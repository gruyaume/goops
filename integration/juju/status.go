package juju

import (
	"encoding/json"
	"fmt"
)

const (
	statusCommand = "status"
)

type ApplicationStatus struct {
	Current string `json:"current"`
	Message string `json:"message"`
	Since   string `json:"since"`
}

type ApplicationStatusOutput struct {
	Address           string            `json:"address"`
	ApplicationStatus ApplicationStatus `json:"application-status"`
}

type StatusOutput struct {
	Applications map[string]ApplicationStatusOutput `json:"applications"`
}

func (j Client) Status() (*StatusOutput, error) {
	args := []string{statusCommand, "--format=json"}

	output, err := j.Runner.Run(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	var statusOutput StatusOutput

	err = json.Unmarshal(output, &statusOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal models: %w", err)
	}

	return &statusOutput, nil
}
