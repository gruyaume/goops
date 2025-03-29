package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

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
