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

// GetRelationIDs retrieves the IDs of all relations for a given endpoint.
// The output is useful as input to:
// - ListRelationUnits
// - GetAppRelationData
// - GetUnitRelationData
// - SetUnitRelationData
// - SetAppRelationData
// - GetRelationModel
func GetRelationIDs(name string) ([]string, error) {
	commandRunner := GetCommandRunner()

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

// GetUnitRelationData retrieves the relation data for a specific unit in a relation by its ID.
// unitID can either be:
// - The remote unit ID which can be retrieved via goops.ListRelationUnits()
// - The local unit ID which you can retrieve via goops.ReadEnv()
func GetUnitRelationData(id string, unitID string) (map[string]string, error) {
	commandRunner := GetCommandRunner()

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

// GetUnitRelationData retrieves the relation data for a specific app in a relation by its ID.
// unitID can either be:
// - The remote unit ID which can be retrieved via goops.ListRelationUnits()
// - The local unit ID which you can retrieve via goops.ReadEnv()
func GetAppRelationData(id string, unitID string) (map[string]string, error) {
	commandRunner := GetCommandRunner()

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

// ListRelationUnits lists all remote units in a relation by its ID.
func ListRelationUnits(id string) ([]string, error) {
	commandRunner := GetCommandRunner()

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

// GetRelationApp retrieves the remote application name for a relation by its ID.
func GetRelationApp(id string) (string, error) {
	commandRunner := GetCommandRunner()

	args := []string{"-r=" + id, "--app", "--format=json"}

	output, err := commandRunner.Run(relationListCommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to list relation data: %w", err)
	}

	var relationApp string

	err = json.Unmarshal(output, &relationApp)
	if err != nil {
		return "", fmt.Errorf("failed to parse relation list: %w", err)
	}

	return relationApp, nil
}

// SetUnitRelationData sets the local unit relation data in a relation by its ID.
func SetUnitRelationData(id string, data map[string]string) error {
	commandRunner := GetCommandRunner()

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

// SetAppRelationData sets the local application relation data in a relation by its ID.
func SetAppRelationData(id string, data map[string]string) error {
	commandRunner := GetCommandRunner()

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

// GetRelationModel retrieves the relation model UUID for a relation by its ID.
func GetRelationModelUUID(id string) (string, error) {
	commandRunner := GetCommandRunner()

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
