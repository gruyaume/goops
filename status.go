package goops

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

type Status struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

func SetUnitStatus(status StatusName, message ...string) error {
	commandRunner := GetCommandRunner()

	var args []string

	args = append(args, string(status))

	if len(message) > 0 {
		args = append(args, message...)
	}

	_, err := commandRunner.Run(statusSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

func SetAppStatus(status StatusName, message ...string) error {
	commandRunner := GetCommandRunner()

	var args []string

	args = append(args, "--application", string(status))

	if len(message) > 0 {
		args = append(args, message...)
	}

	_, err := commandRunner.Run(statusSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

func GetUnitStatus() (*Status, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--include-data", "--format=json"}

	output, err := commandRunner.Run(statusGetCommand, args...)
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

func GetAppStatus() (*Status, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--application", "--include-data", "--format=json"}

	output, err := commandRunner.Run(statusGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get application status: %w", err)
	}

	var status Status

	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse application status: %w", err)
	}

	return &status, nil
}
