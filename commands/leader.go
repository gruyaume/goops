package commands

import (
	"encoding/json"
	"fmt"
)

const (
	isLeaderCommand = "is-leader"
)

func (command Command) IsLeader() (bool, error) {
	args := []string{"--format=json"}

	output, err := command.Runner.Run(isLeaderCommand, args...)
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
