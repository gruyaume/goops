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

// AddStorage adds a storage instance to the unit.
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

type StorageInfo struct {
	Kind     string `json:"kind"`
	Location string `json:"location"`
}

// GetStorageByID retrieves storage information by its ID.
func GetStorageByID(id string) (*StorageInfo, error) {
	commandRunner := GetCommandRunner()

	args := []string{"-s", id}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(storageGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %w", err)
	}

	var storageInfo StorageInfo

	err = json.Unmarshal(output, &storageInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse storage: %w", err)
	}

	return &storageInfo, nil
}

// ListStorage lists all storage IDs for a given storage name.
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
