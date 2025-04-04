package commands

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	secretAddCommand   = "secret-add"
	secretGetCommand   = "secret-get"
	secretGrantCommand = "secret-grant"
	secredIDsCommand   = "secret-ids"
	secretInfoGet      = "secret-info-get"
	secretRemove       = "secret-remove"
	secretRevoke       = "secret-revoke"
	secretSet          = "secret-set"
)

type SecretInfo struct {
	Revision int    `json:"revision"`
	Label    string `json:"label"`
	Owner    string `json:"owner"`
	Rotation string `json:"rotation"`
}

func (command Command) SecretAdd(content map[string]string, description string, expire time.Time, label string, owner string, rotate string) (string, error) {
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

	if owner != "" {
		args = append(args, "--owner="+owner)
	}

	if rotate != "" {
		args = append(args, "--rotate="+rotate)
	}

	if !expire.IsZero() {
		args = append(args, "--expire="+expire.Format(time.RFC3339))
	}

	output, err := command.Runner.Run(secretAddCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to add secret: %w", err)
	}

	return string(output), nil
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

func (command Command) SecretGrant(id string, relation string, unit string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	if relation == "" {
		return fmt.Errorf("relation cannot be empty")
	}

	args := []string{id, "--relation=" + relation}

	if unit != "" {
		args = append(args, "--unit="+unit)
	}

	_, err := command.Runner.Run(secretGrantCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to grant secret: %w", err)
	}

	return nil
}

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

func (command Command) SecretInfoGet(id string, label string) (map[string]SecretInfo, error) {
	if id == "" && label == "" {
		return nil, fmt.Errorf("either secret ID or label must be provided")
	}

	if id != "" && label != "" {
		return nil, fmt.Errorf("only one of secret ID or label can be provided")
	}

	args := []string{}
	if id != "" {
		args = append(args, id)
	}

	if label != "" {
		args = append(args, "--label="+label)
	}

	args = append(args, "--format=json")

	output, err := command.Runner.Run(secretInfoGet, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret info: %w", err)
	}

	var secretInfo map[string]SecretInfo

	err = json.Unmarshal(output, &secretInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret info: %w", err)
	}

	if len(secretInfo) == 0 {
		return nil, fmt.Errorf("no secret info found for ID or label: %s", id)
	}

	return secretInfo, nil
}

func (command Command) SecretRemove(id string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	args := []string{id}

	_, err := command.Runner.Run(secretRemove, args...)
	if err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	return nil
}

func (command Command) SecretRevoke(id string, unit string, app string, relation string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	args := []string{id}

	if unit != "" {
		args = append(args, "--unit="+unit)
	}

	if app != "" {
		args = append(args, "--app="+app)
	}

	if relation != "" {
		args = append(args, "--relation="+relation)
	}

	_, err := command.Runner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

func (command Command) SecretSet(id string, content map[string]string, description string, expire time.Time, label string, owner string, rotate string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	if len(content) == 0 {
		return fmt.Errorf("content cannot be empty")
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

	if owner != "" {
		args = append(args, "--owner="+owner)
	}

	if rotate != "" {
		args = append(args, "--rotate="+rotate)
	}

	if !expire.IsZero() {
		args = append(args, "--expire="+expire.Format(time.RFC3339))
	}

	_, err := command.Runner.Run(secretSet, args...)
	if err != nil {
		return fmt.Errorf("failed to set secret: %w", err)
	}

	return nil
}
