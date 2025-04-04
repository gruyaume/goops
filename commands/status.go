package commands

import (
	"encoding/json"
	"fmt"
)

type StatusName string

const (
	StatusActive      StatusName = "active"
	StatusBlocked     StatusName = "blocked"
	StatusWaiting     StatusName = "waiting"
	StatusMaintenance StatusName = "maintenance"
)

const (
	statusGetCommand = "status-get"
	statusSetCommand = "status-set"
)

type StatusOptions struct {
	Name    StatusName
	Message string
}

type Status struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

func (command Command) StatusSet(opts *StatusOptions) error {
	var args []string

	args = append(args, string(opts.Name))
	if opts.Message != "" {
		args = append(args, opts.Message)
	}

	_, err := command.Runner.Run(statusSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

func (command Command) StatusGet() (*Status, error) {
	args := []string{"--include-data", "--format=json"}

	output, err := command.Runner.Run(statusGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var status Status

	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}

	return &status, nil
}
