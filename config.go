package goops

import (
	"encoding/json"
	"fmt"
)

const (
	configGetCommand = "config-get"
)

// GetConfig retrieves the Juju configuration options and unmarshals them into the provided config struct.
func GetConfig(config any) error {
	commandRunner := GetCommandRunner()

	args := []string{"--all", "--format=json"}

	output, err := commandRunner.Run(configGetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	err = json.Unmarshal(output, config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
