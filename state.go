package goops

import (
	"encoding/json"
	"fmt"
)

const (
	stateDeleteCommand = "state-delete"
	stateGetCommand    = "state-get"
	stateSetCommand    = "state-set"
)

func DeleteState(key string) error {
	commandRunner := GetCommandRunner()

	args := []string{key}

	output, err := commandRunner.Run(stateDeleteCommand, args...)
	if err != nil {
		return err
	}

	if len(output) > 0 {
		return fmt.Errorf("unexpected output: %s", output)
	}

	return nil
}

func GetState(key string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{key, "--format=json"}

	output, err := commandRunner.Run(stateGetCommand, args...)
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

func SetState(key string, value string) error {
	commandRunner := GetCommandRunner()

	args := []string{key + "=" + value}

	_, err := commandRunner.Run(stateSetCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
