package goops

import "encoding/json"

const (
	goalStateCommand = "goal-state"
)

type GoalStateStatusContents struct {
	Status string `json:"status"`
	Since  string `json:"since,omitempty"`
}

type UnitsGoalStateContents map[string]GoalStateStatusContents

type GoalState struct {
	Units     UnitsGoalStateContents            `json:"units"`
	Relations map[string]UnitsGoalStateContents `json:"relations"`
}

// GetGoalState retrieves the status of the charm's peers and related units.
func GetGoalState() (*GoalState, error) {
	commandRunner := GetCommandRunner()

	args := []string{"--format=json"}

	output, err := commandRunner.Run(goalStateCommand, args...)
	if err != nil {
		return nil, err
	}

	var goalState GoalState

	err = json.Unmarshal(output, &goalState)
	if err != nil {
		return nil, err
	}

	if len(goalState.Relations) == 0 {
		goalState.Relations = nil
	}

	return &goalState, nil
}
