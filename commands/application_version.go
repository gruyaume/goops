package commands

import (
	"fmt"
)

const (
	applicationVersionSetCommand = "application-version-set"
)

type ApplicationVersionSetOptions struct {
	Version string
}

func (command Command) ApplicationVersionSet(opts *ApplicationVersionSetOptions) error {
	args := []string{}
	if opts.Version != "" {
		args = append(args, opts.Version)
	}

	_, err := command.Runner.Run(applicationVersionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set application version: %w", err)
	}

	return nil
}
