package commands

import (
	"fmt"
)

const (
	applicationVersionSetCommand = "application-version-set"
)

func (command Command) ApplicationVersionSet(message string) error {
	args := []string{}
	if message != "" {
		args = append(args, message)
	}

	_, err := command.Runner.Run(applicationVersionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set application version: %w", err)
	}

	return nil
}
