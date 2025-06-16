package goops

import "encoding/json"

const (
	goalStateCommand = "goal-state"
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
