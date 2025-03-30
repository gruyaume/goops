package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CommandRunner is an interface for running commands.
// It allows for mocking in tests.
type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

// HookCommand is the default implementation of CommandRunner using os/exec.
type HookCommand struct{}

func (r *HookCommand) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command %s failed: %s: %w", name, stderr.String(), err)
	}
	return output, nil
}

type Command struct {
	Runner CommandRunner
}

func NewCommand() *Command {
	return &Command{
		Runner: &HookCommand{},
	}
}
