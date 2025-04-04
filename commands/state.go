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

type StateDeleteOptions struct {
	Key string
}

type StateGetOptions struct {
	Key string
}

type StateSetOptions struct {
	Key   string
	Value string
}

func (command Command) StateDelete(opts *StateDeleteOptions) error {
	args := []string{opts.Key}

	output, err := command.Runner.Run(stateDeleteCommand, args...)
	if err != nil {
		return err
	}

	if len(output) > 0 {
		return fmt.Errorf("unexpected output: %s", output)
	}

	return nil
}

func (command Command) StateGet(opts *StateGetOptions) (string, error) {
	if opts.Key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	args := []string{opts.Key, "--format=json"}

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
		return "", fmt.Errorf("no state found for key: %s", opts.Key)
	}

	return state, nil
}

func (command Command) StateSet(opts *StateSetOptions) error {
	if opts.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if opts.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	args := []string{opts.Key + "=" + opts.Value}

	_, err := command.Runner.Run(stateSetCommand, args...)
	if err != nil {
		return err
	}

	return nil
}
