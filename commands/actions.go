package commands

import (
	"encoding/json"
	"fmt"
)

const (
	actionFailCommand = "action-fail"
	actionGetCommand  = "action-get"
	actionLogCommand  = "action-log"
	actionSetCommand  = "action-set"
)

type ActionFailOptions struct {
	Message string
}

type ActionGetOptions struct {
	Key string
}

type ActionLogOptions struct {
	Message string
}

type ActionSetOptions struct {
	Content map[string]string
}

func (command Command) ActionFail(opts *ActionFailOptions) error {
	args := []string{}
	if opts.Message != "" {
		args = append(args, opts.Message)
	}

	_, err := command.Runner.Run(actionFailCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to fail action: %w", err)
	}

	return nil
}

func (command Command) ActionGet(opts *ActionGetOptions) (string, error) {
	args := []string{opts.Key, "--format=json"}

	output, err := command.Runner.Run(actionGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get action parameter: %w", err)
	}

	var actionParameter string

	err = json.Unmarshal(output, &actionParameter)
	if err != nil {
		return "", fmt.Errorf("failed to parse action parameter: %w", err)
	}

	return actionParameter, nil
}

func (command Command) ActionLog(opts *ActionLogOptions) error {
	if opts.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	args := []string{opts.Message}

	_, err := command.Runner.Run(actionLogCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to log action message: %w", err)
	}

	return nil
}

func (command Command) ActionSet(opts *ActionSetOptions) error {
	if opts.Content == nil {
		return fmt.Errorf("content cannot be empty")
	}

	var args []string
	for key, value := range opts.Content {
		args = append(args, key+"="+value)
	}

	_, err := command.Runner.Run(actionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set action parameters: %w", err)
	}

	return nil
}
