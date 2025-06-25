package goops

import (
	"encoding/json"
	"fmt"
)

const (
	isLeaderCommand = "is-leader"
)

// IsLeader retrieves the unit's leadership status.
func IsLeader() (bool, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--format=json"}

	output, err := commandRunner.Run(isLeaderCommand, args...)
	if err != nil {
		return false, fmt.Errorf("failed to verify if unit is leader: %w", err)
	}

	var isLeader bool

	err = json.Unmarshal(output, &isLeader)
	if err != nil {
		return false, fmt.Errorf("failed to parse leader status: %w", err)
	}

	return isLeader, nil
}
