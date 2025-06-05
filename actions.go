package goops

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

func FailActionf(format string, args ...any) error {
	commandRunner := GetRunner()

	message := fmt.Sprintf(format, args...)

	_, err := commandRunner.Run(actionFailCommand, message)
	if err != nil {
		return fmt.Errorf("failed to fail action: %w", err)
	}

	return nil
}

func GetActionParameter(key string) (string, error) {
	commandRunner := GetRunner()

	args := []string{key, "--format=json"}

	output, err := commandRunner.Run(actionGetCommand, args...)
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

func LogActionf(format string, args ...any) error {
	commandRunner := GetRunner()

	message := fmt.Sprintf(format, args...)

	_, err := commandRunner.Run(actionLogCommand, message)
	if err != nil {
		return fmt.Errorf("failed to log action message: %w", err)
	}

	return nil
}

func SetActionResults(results map[string]string) error {
	commandRunner := GetRunner()

	var args []string

	for key, value := range results {
		args = append(args, key+"="+value)
	}

	_, err := commandRunner.Run(actionSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set action parameters: %w", err)
	}

	return nil
}
