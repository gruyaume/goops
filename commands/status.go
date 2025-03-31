package commands

import (
	"fmt"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
	StatusWaiting Status = "waiting"
)

const (
	StatusSetCommand = "status-set"
)

func (command Command) StatusSet(status Status, message string) error {
	var args []string

	args = append(args, string(status))
	if message != "" {
		args = append(args, message)
	}

	_, err := command.Runner.Run(StatusSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}
