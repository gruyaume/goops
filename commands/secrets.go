package commands

import (
	"encoding/json"
	"fmt"
)

const (
	secredIDsCommand = "secret-ids"
	secretGetCommand = "secret-get"
	secretAddCommand = "secret-add"
)

func (command Command) SecretIDs() ([]string, error) {
	output, err := command.Runner.Run(secredIDsCommand, "--format=json")
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

func (command Command) SecretGet(id string, label string, peek bool, refresh bool) (map[string]string, error) {
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

	args = append(args, "--format=json")

	output, err := command.Runner.Run(secretGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	var secretContent map[string]string

	err = json.Unmarshal(output, &secretContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret content: %w", err)
	}

	return secretContent, nil
}

func (command Command) SecretAdd(content map[string]string, description string, label string) (string, error) {
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

	output, err := command.Runner.Run(secretAddCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to add secret: %w", err)
	}

	return string(output), nil
}
