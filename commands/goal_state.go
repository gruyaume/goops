package commands

import "encoding/json"

const (
	GoalStateCommand = "goal-state"
)

type UnitStatus struct {
	Status string `json:"status"`
	Since  string `json:"since"`
}

type RelationStatus struct {
	Status string `json:"status"`
	Since  string `json:"since"`
}

type GoalState struct {
	Units     map[string]*UnitStatus                `json:"units"`
	Relations map[string]map[string]*RelationStatus `json:"relations"`
}

func (command Command) GoalState() (*GoalState, error) {
	args := []string{"--format=json"}

	output, err := command.Runner.Run(GoalStateCommand, args...)
	if err != nil {
		return nil, err
	}

	// command.JujuLog(Warning, "goal-state output: %s", string(output))

	var goalState GoalState

	err = json.Unmarshal(output, &goalState)
	if err != nil {
		return nil, err
	}

	return &goalState, nil
}
