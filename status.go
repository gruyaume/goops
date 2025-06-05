package goops

import "fmt"

type Status string

const (
	StatusActive      Status = "active"
	StatusBlocked     Status = "blocked"
	StatusWaiting     Status = "waiting"
	StatusMaintenance Status = "maintenance"
)

const (
	statusSetCommand = "status-set"
)

func SetUnitStatus(status Status, message ...string) error {
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
