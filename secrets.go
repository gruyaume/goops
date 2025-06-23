package goops

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

type SecretRotate string

const (
	RotateHourly  SecretRotate = "hourly"
	RotateDaily   SecretRotate = "daily"
	RotateMonthly SecretRotate = "monthly"
	RotateNever   SecretRotate = "never"
)

type SetSecretOptions struct {
	ID          string
	Content     map[string]string
	Description string
	Expire      time.Time
	Label       string
	Owner       SecretOwner
	Rotate      SecretRotate
}

type SecretOwner string

const (
	OwnerApplication SecretOwner = "application"
	OwnerUnit        SecretOwner = "unit"
)

type AddSecretOptions struct {
	Content     map[string]string
	Description string
	Expire      time.Time
	Label       string
	Owner       SecretOwner
	Rotate      SecretRotate
}

// AddSecret adds a new secret with the provided options.
func AddSecret(opts *AddSecretOptions) (string, error) {
	commandRunner := GetCommandRunner()

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
		args = append(args, "--owner="+string(opts.Owner))
	}

	if opts.Rotate != "" {
		args = append(args, "--rotate="+string(opts.Rotate))
	}

	if !opts.Expire.IsZero() {
		args = append(args, "--expire="+opts.Expire.Format(time.RFC3339))
	}

	output, err := commandRunner.Run(secretAddCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to add secret: %w", err)
	}

	return string(output), nil
}

// GetSecretByID retrieves the secret content by its ID.
func GetSecretByID(id string, peek bool, refresh bool) (map[string]string, error) {
	commandRunner := GetCommandRunner()

	var args []string
	args = append(args, id)

	if peek {
		args = append(args, "--peek")
	}

	if refresh {
		args = append(args, "--refresh")
	}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(secretGetCommand, args...)
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

// GetSecretByLabel retrieves the secret content by its label.
func GetSecretByLabel(label string, peek bool, refresh bool) (map[string]string, error) {
	commandRunner := GetCommandRunner()

	var args []string

	args = append(args, "--label="+label)

	if peek {
		args = append(args, "--peek")
	}

	if refresh {
		args = append(args, "--refresh")
	}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(secretGetCommand, args...)
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

// GrantSecretToRelation grants a secret to a specific relation.
// All units of the related application are granted access
func GrantSecretToRelation(id string, relation string) error {
	commandRunner := GetCommandRunner()

	args := []string{id, "--relation=" + relation}

	_, err := commandRunner.Run(secretGrantCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to grant secret: %w", err)
	}

	return nil
}

// GrantSecretToUnit grants a secret to a specific unit in a relation.
func GrantSecretToUnit(id string, relation string, unit string) error {
	commandRunner := GetCommandRunner()

	args := []string{id, "--relation=" + relation, "--unit=" + unit}

	_, err := commandRunner.Run(secretGrantCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to grant secret: %w", err)
	}

	return nil
}

// GetSecretIDs retrieves the IDs for secrets owned by the application.
func GetSecretIDs() ([]string, error) {
	commandRunner := GetCommandRunner()

	output, err := commandRunner.Run(secredIDsCommand, "--format=json")
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

// GetSecretInfoByID retrieves a secret metadata info by its ID.
func GetSecretInfoByID(id string) (map[string]SecretInfo, error) {
	commandRunner := GetCommandRunner()

	args := []string{}

	args = append(args, id)

	args = append(args, "--format=json")

	output, err := commandRunner.Run(secretInfoGet, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret info: %w", err)
	}

	var secretInfo map[string]SecretInfo

	err = json.Unmarshal(output, &secretInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret info: %w", err)
	}

	if len(secretInfo) == 0 {
		return nil, fmt.Errorf("no secret info found for ID: %s", id)
	}

	return secretInfo, nil
}

// GetSecretInfoByLabel retrieves a secret metadata info by its label.
func GetSecretInfoByLabel(label string) (map[string]SecretInfo, error) {
	commandRunner := GetCommandRunner()

	args := []string{}

	args = append(args, "--label="+label)

	args = append(args, "--format=json")

	output, err := commandRunner.Run(secretInfoGet, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret info: %w", err)
	}

	var secretInfo map[string]SecretInfo

	err = json.Unmarshal(output, &secretInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret info: %w", err)
	}

	if len(secretInfo) == 0 {
		return nil, fmt.Errorf("no secret info found for label: %s", label)
	}

	return secretInfo, nil
}

// RemoveSecret removes a secret by its ID.
func RemoveSecret(id string) error {
	commandRunner := GetCommandRunner()

	args := []string{id}

	_, err := commandRunner.Run(secretRemove, args...)
	if err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	return nil
}

// RevokeSecret revokes a secret by its ID.
func RevokeSecret(id string) error {
	commandRunner := GetCommandRunner()

	args := []string{id}

	_, err := commandRunner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

// RevokeSecretFromRelation revokes a secret from a specific relation.
func RevokeSecretFromRelation(id string, relation string) error {
	commandRunner := GetCommandRunner()

	args := []string{id}

	args = append(args, "--relation="+relation)

	_, err := commandRunner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

// RevokeSecretFromApp revokes a secret from a specific application.
func RevokeSecretFromApp(id string, app string) error {
	commandRunner := GetCommandRunner()

	args := []string{id}

	args = append(args, "--app="+app)

	_, err := commandRunner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

// RevokeSecretFromApp revokes a secret from a specific application.
func RevokeSecretFromUnit(id string, unit string) error {
	commandRunner := GetCommandRunner()

	args := []string{id}

	args = append(args, "--unit="+unit)

	_, err := commandRunner.Run(secretRevoke, args...)
	if err != nil {
		return fmt.Errorf("failed to revoke secret: %w", err)
	}

	return nil
}

// SetSecret updates an existing secret with new content and options.
func SetSecret(opts *SetSecretOptions) error {
	commandRunner := GetCommandRunner()

	if opts.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}

	args := []string{opts.ID}

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
		args = append(args, "--owner="+string(opts.Owner))
	}

	if opts.Rotate != "" {
		args = append(args, "--rotate="+string(opts.Rotate))
	}

	if !opts.Expire.IsZero() {
		args = append(args, "--expire="+opts.Expire.Format(time.RFC3339))
	}

	_, err := commandRunner.Run(secretSet, args...)
	if err != nil {
		return fmt.Errorf("failed to set secret: %w", err)
	}

	return nil
}
