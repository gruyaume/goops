package goops

import (
	"bytes"
	"fmt"
	"os/exec"
)

var defaultRunner CommandRunner = &realHookCommand{}

type realHookCommand struct{}

func (r *realHookCommand) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command %s failed: %s: %w", name, stderr.String(), err)
	}

	return output, nil
}

// CommandRunner is an interface for running commands.
// It allows for mocking in tests.
type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

func GetCommandRunner() CommandRunner {
	return defaultRunner
}

func SetCommandRunner(runner CommandRunner) {
	defaultRunner = runner
}
