package goops

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

func GetRelationIDs(name string) ([]string, error) {
	commandRunner := GetRunner()

	args := []string{name, "--format=json"}

	output, err := commandRunner.Run(relationIDsCommand, args...)
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

func GetUnitRelationData(id string, unitID string) (map[string]string, error) {
	commandRunner := GetRunner()

	args := []string{"-r=" + id, "-", unitID}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(relationGetCommand, args...)
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

func GetAppRelationData(id string, unitID string) (map[string]string, error) {
	commandRunner := GetRunner()

	args := []string{"-r=" + id, "-", unitID, "--app"}

	args = append(args, "--format=json")

	output, err := commandRunner.Run(relationGetCommand, args...)
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

func ListRelations(id string) ([]string, error) {
	commandRunner := GetRunner()

	args := []string{"-r=" + id, "--format=json"}

	output, err := commandRunner.Run(relationListCommand, args...)
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

func SetUnitRelationData(id string, data map[string]string) error {
	commandRunner := GetRunner()

	args := []string{"-r=" + id}

	for key, value := range data {
		args = append(args, key+"="+value)
	}

	output, err := commandRunner.Run(relationSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set relation data: %w", err)
	}

	if len(output) > 0 {
		return fmt.Errorf("failed to set relation data: %s", output)
	}

	return nil
}

func SetAppRelationData(id string, data map[string]string) error {
	commandRunner := GetRunner()

	args := []string{"-r=" + id}

	args = append(args, "--app")

	for key, value := range data {
		args = append(args, key+"="+value)
	}

	output, err := commandRunner.Run(relationSetCommand, args...)
	if err != nil {
		return fmt.Errorf("failed to set relation data: %w", err)
	}

	if len(output) > 0 {
		return fmt.Errorf("failed to set relation data: %s", output)
	}

	return nil
}

type RelationModel struct {
	UUID string `json:"uuid"`
}

func GetRelationModel(id string) (string, error) {
	commandRunner := GetRunner()

	args := []string{"-r=" + id, "--format=json"}

	output, err := commandRunner.Run(relationModelGetCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get relation model data: %w", err)
	}

	var relationModel RelationModel

	err = json.Unmarshal(output, &relationModel)
	if err != nil {
		return "", fmt.Errorf("failed to parse relation model data: %w", err)
	}

	return relationModel.UUID, nil
}
