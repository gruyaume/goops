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

func StatusSet(runner CommandRunner, status Status) error {
	_, err := runner.Run(StatusSetCommand, string(status))
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}
	return nil
}
