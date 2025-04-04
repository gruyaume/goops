package commands

import (
	"encoding/json"
	"fmt"
)

const (
	stateDeleteCommand = "state-delete"
	stateGetCommand    = "state-get"
	stateSetCommand    = "state-set"
)

func (command Command) StateDelete(key string) error {
	args := []string{key}

	output, err := command.Runner.Run(stateDeleteCommand, args...)
	if err != nil {
		return err
	}

	if len(output) > 0 {
		return fmt.Errorf("unexpected output: %s", output)
	}

	return nil
}

func (command Command) StateGet(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	args := []string{key, "--format=json"}

	output, err := command.Runner.Run(stateGetCommand, args...)
	if err != nil {
		return "", err
	}

	var state string

	err = json.Unmarshal(output, &state)
	if err != nil {
		return "", fmt.Errorf("failed to parse state: %w", err)
	}

	if len(state) == 0 {
		return "", fmt.Errorf("no state found for key: %s", key)
	}

	return state, nil
}

func (command Command) StateSet(key string, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	args := []string{key + "=" + value}

	_, err := command.Runner.Run(stateSetCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
