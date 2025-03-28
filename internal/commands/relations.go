package commands

import (
	"encoding/json"
	"fmt"
)

const (
	RelationIDsCommand = "relation-ids"
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
