package commands

import (
	"encoding/json"
	"fmt"
)

const (
	ActionGetCommand  = "action-get"
	ActionFailCommand = "action-fail"
	ActionSetCommand  = "action-set"
)

func ActionGet(runner CommandRunner, key string) (string, error) {
	args := []string{key, "--format=json"}
	output, err := runner.Run(ActionGetCommand, args...)
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

func ActionFail(runner CommandRunner, message string) error {
	args := []string{}
	if message != "" {
		args = append(args, message)
	}
	_, err := runner.Run(ActionFailCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to fail action: %w", err)
	}

	return nil
}

func ActionSet(runner CommandRunner, content map[string]string) error {
	if content == nil {
		return fmt.Errorf("content cannot be empty")
	}
	var args []string
	for key, value := range content {
		args = append(args, key+"="+value)
	}
	_, err := runner.Run(ActionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set action parameters: %w", err)
	}
	return nil
}
