package commands

import (
	"fmt"
	"os/exec"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
	StatusWaiting Status = "waiting"
)

const (
	SetStatusCommand = "status-set"
)

func SetStatus(status Status) error {
	cmd := exec.Command(SetStatusCommand, string(status))
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}
	return nil
}
