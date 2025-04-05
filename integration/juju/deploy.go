package juju

import (
	"fmt"
)

const (
	deployCommand = "deploy"
)

type DeployOptions struct {
	Charm           string
	ApplicationName string
}

func (j Client) Deploy(opts *DeployOptions) error {
	if opts.Charm == "" {
		return fmt.Errorf("charm is required")
	}

	args := []string{deployCommand, opts.Charm}
	if opts.ApplicationName != "" {
		args = append(args, opts.ApplicationName)
	}

	_, err := j.Runner.Run(args...)
	if err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}
