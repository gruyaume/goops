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

func (command Command) ActionGet(key string) (string, error) {
	args := []string{key, "--format=json"}
	output, err := command.Runner.Run(ActionGetCommand, args...)
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

func (command Command) ActionFail(message string) error {
	args := []string{}
	if message != "" {
		args = append(args, message)
	}
	_, err := command.Runner.Run(ActionFailCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to fail action: %w", err)
	}

	return nil
}

func (command Command) ActionSet(content map[string]string) error {
	if content == nil {
		return fmt.Errorf("content cannot be empty")
	}
	var args []string
	for key, value := range content {
		args = append(args, key+"="+value)
	}
	_, err := command.Runner.Run(ActionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set action parameters: %w", err)
	}
	return nil
}
