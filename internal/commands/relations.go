package commands

import (
	"encoding/json"
	"fmt"
)

const (
	RelationIDsCommand  = "relation-ids"
	RelationGetCommand  = "relation-get"
	RelationListCommand = "relation-list"
)

func RelationIDs(runner CommandRunner, name string) ([]string, error) {
	args := []string{name, "--format=json"}
	output, err := runner.Run(RelationIDsCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get relation IDs: %w", err)
	}
	var relationIDs []string
	err = json.Unmarshal(output, &relationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation IDs: %w", err)
	}
	return relationIDs, nil
}

func RelationGet(runner CommandRunner, id string, unitID string, app bool) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("relation ID is empty")
	}
	if unitID == "" {
		return nil, fmt.Errorf("unit ID is empty")
	}
	args := []string{"-r=" + id, "-", unitID}
	if app {
		args = append(args, "--app")
	}
	args = append(args, "--format=json")
	output, err := runner.Run(RelationGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get relation data: %w", err)
	}
	var relationContent map[string]string
	err = json.Unmarshal(output, &relationContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation content: %w", err)
	}
	return relationContent, nil
}

func RelationList(runner CommandRunner, id string) ([]string, error) {
	if id == "" {
		return nil, fmt.Errorf("relation ID is empty")
	}
	args := []string{"-r=" + id, "--format=json"}
	output, err := runner.Run(RelationListCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list relation data: %w", err)
	}
	var relationList []string
	err = json.Unmarshal(output, &relationList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation list: %w", err)
	}
	return relationList, nil
}
