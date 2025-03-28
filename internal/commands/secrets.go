package commands

import (
	"encoding/json"
	"fmt"
)

const (
	SecredIDsCommand = "secret-ids"
	SecretGetCommand = "secret-get"
	SecretAddCommand = "secret-add"
)

func SecretIDs(runner CommandRunner) ([]string, error) {
	output, err := runner.Run(SecredIDsCommand, "--format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to get secret IDs: %w", err)
	}
	var secretIDs []string
	err = json.Unmarshal(output, &secretIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret IDs: %w", err)
	}
	return secretIDs, nil
}

func SecretGet(runner CommandRunner, id string, label string, peek bool, refresh bool) (string, error) {
	var args []string
	if id != "" {
		args = append(args, id)
	}
	if label != "" {
		args = append(args, "--label="+label)
	}
	if peek {
		args = append(args, "--peek")
	}
	if refresh {
		args = append(args, "--refresh")
	}
	output, err := runner.Run(SecretGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}
	return string(output), nil
}

func SecretAdd(runner CommandRunner, content map[string]string, description string, label string) (string, error) {

	if len(content) == 0 {
		return "", fmt.Errorf("content cannot be empty")
	}
	var args []string
	for key, value := range content {
		args = append(args, key+"="+value)
	}
	if description != "" {
		args = append(args, "--description="+description)
	}
	if label != "" {
		args = append(args, "--label="+label)
	}
	output, err := runner.Run(SecretAddCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to add secret: %w", err)
	}
	return string(output), nil
}
