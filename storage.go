package goops

import (
	"encoding/json"
	"fmt"
)

const (
	storageAddCommand  = "storage-add"
	storageGetCommand  = "storage-get"
	storageListCommand = "storage-list"
)

func AddStorage(name string, count int) error {
	commandRunner := GetCommandRunner()

	args := []string{}

	args = append(args, name+"="+fmt.Sprintf("%d", count))

	_, err := commandRunner.Run(storageAddCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to add storage: %w", err)
	}

	return nil
}

func GetStorageByID(id string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{id}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(storageGetCommand, args...)
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

func GetStorageByName(name string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{"-s", name}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(storageGetCommand, args...)
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

func ListStorage(name string) ([]string, error) {
	commandRunner := GetCommandRunner()

	args := []string{name, "--format=json"}

	output, err := commandRunner.Run(storageListCommand, args...)
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
