package commands

import (
	"encoding/json"
	"fmt"
)

const (
	IsLeaderCommand = "is-leader"
)

func IsLeader(runner CommandRunner) (bool, error) {
	args := []string{"--format=json"}
	output, err := runner.Run(IsLeaderCommand, args...)
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
