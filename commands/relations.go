package commands

import (
	"encoding/json"
	"fmt"
)

const (
	relationIDsCommand      = "relation-ids"
	relationGetCommand      = "relation-get"
	relationListCommand     = "relation-list"
	relationSetCommand      = "relation-set"
	relationModelGetCommand = "relation-model-get"
)

type RelationModel struct {
	UUID string `json:"uuid"`
}

func (command Command) RelationIDs(name string) ([]string, error) {
	args := []string{name, "--format=json"}

	output, err := command.Runner.Run(relationIDsCommand, args...)
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

func (command Command) RelationGet(id string, unitID string, app bool) (map[string]string, error) {
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

	output, err := command.Runner.Run(relationGetCommand, args...)
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

func (command Command) RelationList(id string) ([]string, error) {
	if id == "" {
		return nil, fmt.Errorf("relation ID is empty")
	}

	args := []string{"-r=" + id, "--format=json"}

	output, err := command.Runner.Run(relationListCommand, args...)
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

func (command Command) RelationSet(id string, app bool, data map[string]string) error {
	if id == "" {
		return fmt.Errorf("relation ID is empty")
	}

	args := []string{"-r=" + id}
	if app {
		args = append(args, "--app")
	}

	for key, value := range data {
		args = append(args, key+"="+value)
	}

	output, err := command.Runner.Run(relationSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set relation data: %w", err)
	}

	if len(output) > 0 {
		return fmt.Errorf("failed to set relation data: %s", output)
	}

	return nil
}

func (command Command) RelationModelGet(id string) (*RelationModel, error) {
	if id == "" {
		return nil, fmt.Errorf("relation ID is empty")
	}

	args := []string{"-r=" + id, "--format=json"}

	output, err := command.Runner.Run(relationModelGetCommand, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get relation model data: %w", err)
	}

	var relationModel RelationModel

	err = json.Unmarshal(output, &relationModel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation model data: %w", err)
	}

	return &relationModel, nil
}
