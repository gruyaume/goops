package goops

import (
	"encoding/json"
	"fmt"
)

type StatusCode string

const (
	StatusActive      StatusCode = "active"
	StatusBlocked     StatusCode = "blocked"
	StatusWaiting     StatusCode = "waiting"
	StatusMaintenance StatusCode = "maintenance"
)

const (
	statusGetCommand = "status-get"
	statusSetCommand = "status-set"
)

type Status struct {
	Code    StatusCode `json:"status"`
	Message string     `json:"message"`
}

func SetUnitStatus(status StatusCode, message ...string) error {
	commandRunner := GetRunner()

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

func SetAppStatus(status StatusCode, message ...string) error {
	commandRunner := GetRunner()

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

func GetStatus() (*Status, error) {
	commandRunner := GetRunner()

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
