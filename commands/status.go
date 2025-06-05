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
)

type Status struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
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
