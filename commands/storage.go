package commands

import (
	"encoding/json"
	"fmt"
)

const (
	storageAddCommand  = "storage-add"
	storageGetCommand  = "storage-get"
	storageListCommand = "storage-list"
)

type StorageAddOptions struct {
	Name  string
	Count int
}

type StorageGetOptions struct {
	ID   string
	Name string
}

type StorageListOptions struct {
	Name string
}

func (command Command) StorageAdd(opts *StorageAddOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("storage name cannot be empty")
	}

	args := []string{}

	if opts.Count > 0 {
		args = append(args, opts.Name+"="+fmt.Sprintf("%d", opts.Count))
	} else {
		args = append(args, opts.Name)
	}

	_, err := command.Runner.Run(storageAddCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to add storage: %w", err)
	}

	return nil
}

func (command Command) StorageGet(opts *StorageGetOptions) (string, error) {
	if opts.ID == "" && opts.Name == "" {
		return "", fmt.Errorf("either ID or Name must be provided")
	}

	args := []string{}

	if opts.ID != "" {
		args = append(args, opts.ID)
	}

	if opts.Name != "" {
		args = append(args, "-s", opts.Name)
	}

	args = append(args, "--format=json")

	output, err := command.Runner.Run(storageGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get storage: %w", err)
	}

	var storage string

	err = json.Unmarshal(output, &storage)
	if err != nil {
		return "", fmt.Errorf("failed to parse storage: %w", err)
	}

	return storage, nil
}

func (command Command) StorageList(opts *StorageListOptions) ([]string, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("storage name cannot be empty")
	}

	args := []string{opts.Name, "--format=json"}

	output, err := command.Runner.Run(storageListCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list storage: %w", err)
	}

	var storageNames []string

	err = json.Unmarshal(output, &storageNames)
	if err != nil {
		return nil, fmt.Errorf("failed to parse storages: %w", err)
	}

	return storageNames, nil
}
