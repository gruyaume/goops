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
	StatusUnknown     StatusName = "unknown"
)

const (
	statusGetCommand = "status-get"
	statusSetCommand = "status-set"
)

type UnitStatus struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

// SetUnitStatus sets the unit status.
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

// SetAppStatus sets the application status.
// Only the leader unit can set the application status.
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

// GetUnitStatus returns the unit status information.
func GetUnitStatus() (*UnitStatus, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--include-data", "--format=json"}

	output, err := commandRunner.Run(statusGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var status UnitStatus

	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}

	return &status, nil
}

type AppStatus struct {
	Name    StatusName `json:"status"`
	Message string     `json:"message"`
}

type appStatusReturn struct {
	AppStatus AppStatus `json:"application-status"`
}

// GetAppStatus returns the application status information.
// Only the leader unit can retrieve the application status.
func GetAppStatus() (*AppStatus, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--application", "--include-data", "--format=json"}

	output, err := commandRunner.Run(statusGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get application status: %w", err)
	}

	var status appStatusReturn

	err = json.Unmarshal(output, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse application status: %w", err)
	}

	return &status.AppStatus, nil
}
