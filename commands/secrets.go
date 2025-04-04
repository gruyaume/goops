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

type SecretAddOptions struct {
	Content     map[string]string
	Description string
	Expire      time.Time
	Label       string
	Owner       string
	Rotate      string // allowed values: hourly, daily, monthly, never
}

type SecretGetOptions struct {
	ID      string
	Label   string
	Peek    bool
	Refresh bool
}

type SecretGrantOptions struct {
	ID       string
	Relation string
	Unit     string
}

type SecretInfoGetOptions struct {
	ID    string
	Label string
}

type SecretRemoveOptions struct {
	ID string
}

type SecretRevokeOptions struct {
	ID       string
	Unit     string
	App      string
	Relation string
}

type SecretSetOptions struct {
	ID          string
	Content     map[string]string
	Description string
	Expire      time.Time
	Label       string
	Owner       string
	Rotate      string // allowed values: hourly, daily, monthly, never
}

func (command Command) SecretAdd(opts *SecretAddOptions) (string, error) {
	if len(opts.Content) == 0 {
		return "", fmt.Errorf("content cannot be empty")
	}

	var args []string
	for key, value := range opts.Content {
		args = append(args, key+"="+value)
	}

	if opts.Description != "" {
		args = append(args, "--description="+opts.Description)
	}

	if opts.Label != "" {
		args = append(args, "--label="+opts.Label)
	}

	if opts.Owner != "" {
		args = append(args, "--owner="+opts.Owner)
	}

	if opts.Rotate != "" {
		args = append(args, "--rotate="+opts.Rotate)
	}

	if !opts.Expire.IsZero() {
		args = append(args, "--expire="+opts.Expire.Format(time.RFC3339))
	}

	output, err := command.Runner.Run(secretAddCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to add secret: %w", err)
	}

	return string(output), nil
}

func (command Command) SecretGet(opts *SecretGetOptions) (map[string]string, error) {
	var args []string
	if opts.ID != "" {
		args = append(args, opts.ID)
	}

	if opts.Label != "" {
		args = append(args, "--label="+opts.Label)
	}

	if opts.Peek {
		args = append(args, "--peek")
	}

	if opts.Refresh {
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

func (command Command) SecretGrant(opts *SecretGrantOptions) error {
	if opts.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	if opts.Relation == "" {
		return fmt.Errorf("relation cannot be empty")
	}

	args := []string{opts.ID, "--relation=" + opts.Relation}

	if opts.Unit != "" {
		args = append(args, "--unit="+opts.Unit)
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

func (command Command) SecretInfoGet(opts *SecretInfoGetOptions) (map[string]SecretInfo, error) {
	if opts.ID == "" && opts.Label == "" {
		return nil, fmt.Errorf("either secret ID or label must be provided")
	}

	if opts.ID != "" && opts.Label != "" {
		return nil, fmt.Errorf("only one of secret ID or label can be provided")
	}

	args := []string{}
	if opts.ID != "" {
		args = append(args, opts.ID)
	}

	if opts.Label != "" {
		args = append(args, "--label="+opts.Label)
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
		return nil, fmt.Errorf("no secret info found for ID or label: %s", opts.ID)
	}

	return secretInfo, nil
}

func (command Command) SecretRemove(opts *SecretRemoveOptions) error {
	if opts.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	args := []string{opts.ID}

	_, err := command.Runner.Run(secretRemove, args...)
	if err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	return nil
}

func (command Command) SecretRevoke(opts *SecretRevokeOptions) error {
	if opts.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	args := []string{opts.ID}

	if opts.Unit != "" {
		args = append(args, "--unit="+opts.Unit)
	}

	if opts.App != "" {
		args = append(args, "--app="+opts.App)
	}

	if opts.Relation != "" {
		args = append(args, "--relation="+opts.Relation)
	}

	_, err := command.Runner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

func (command Command) SecretSet(opts *SecretSetOptions) error {
	if opts.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	if len(opts.Content) == 0 {
		return fmt.Errorf("content cannot be empty")
	}

	var args []string
	for key, value := range opts.Content {
		args = append(args, key+"="+value)
	}

	if opts.Description != "" {
		args = append(args, "--description="+opts.Description)
	}

	if opts.Label != "" {
		args = append(args, "--label="+opts.Label)
	}

	if opts.Owner != "" {
		args = append(args, "--owner="+opts.Owner)
	}

	if opts.Rotate != "" {
		args = append(args, "--rotate="+opts.Rotate)
	}

	if !opts.Expire.IsZero() {
		args = append(args, "--expire="+opts.Expire.Format(time.RFC3339))
	}

	_, err := command.Runner.Run(secretSet, args...)
	if err != nil {
		return fmt.Errorf("failed to set secret: %w", err)
	}

	return nil
}
