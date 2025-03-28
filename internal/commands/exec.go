package commands

import (
	"bytes"
	"log"
	"os/exec"
)

type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

type DefaultRunner struct{}

// Run executes the command using exec.Command, capturing stderr.
func (r *DefaultRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing %s: %s", name, stderr.String())
		return nil, err
	}
	return output, nil
}
