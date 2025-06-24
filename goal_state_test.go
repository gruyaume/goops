package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestGoalState_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"units":{"example/0":{"status":"active","since":"2025-04-03 20:05:33Z"}},"relations":{"certificates":{"tls-certificates-requirer":{"status":"joined","since":"2025-04-03 20:12:11Z"},"tls-certificates-requirer/0":{"status":"active","since":"2025-04-03 20:12:23Z"}}}}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	expectedGoalState := goops.GoalState{
		Units: goops.UnitsGoalStateContents{
			"example/0": {
				Status: "active",
				Since:  "2025-04-03 20:05:33Z",
			},
		},
		Relations: map[string]goops.UnitsGoalStateContents{
			"certificates": {
				"tls-certificates-requirer": {
					Status: "joined",
					Since:  "2025-04-03 20:12:11Z",
				},
				"tls-certificates-requirer/0": {
					Status: "active",
					Since:  "2025-04-03 20:12:23Z",
				},
			},
		},
	}

	goalState, err := goops.GetGoalState()
	if err != nil {
		t.Fatalf("GoalState returned an error: %v", err)
	}

	if len(goalState.Units) != len(expectedGoalState.Units) {
		t.Fatalf("Expected %d units, got %d", len(expectedGoalState.Units), len(goalState.Units))
	}

	for unit, status := range goalState.Units {
		expectedStatus, ok := expectedGoalState.Units[unit]
		if !ok {
			t.Fatalf("Unexpected unit %q", unit)
		}

		if status.Status != expectedStatus.Status {
			t.Errorf("Expected status %q for unit %q, got %q", expectedStatus.Status, unit, status.Status)
		}

		if status.Since != expectedStatus.Since {
			t.Errorf("Expected since %q for unit %q, got %q", expectedStatus.Since, unit, status.Since)
		}
	}

	if len(goalState.Relations) != len(expectedGoalState.Relations) {
		t.Fatalf("Expected %d relations, got %d", len(expectedGoalState.Relations), len(goalState.Relations))
	}

	for relation, units := range goalState.Relations {
		expectedUnits, ok := expectedGoalState.Relations[relation]
		if !ok {
			t.Fatalf("Unexpected relation %q", relation)
		}

		if len(units) != len(expectedUnits) {
			t.Fatalf("Expected %d units for relation %q, got %d", len(expectedUnits), relation, len(units))
		}

		for unit, status := range units {
			expectedStatus, ok := expectedUnits[unit]
			if !ok {
				t.Fatalf("Unexpected unit %q in relation %q", unit, relation)
			}

			if status.Status != expectedStatus.Status {
				t.Errorf("Expected status %q for unit %q in relation %q, got %q", expectedStatus.Status, unit, relation, status.Status)
			}

			if status.Since != expectedStatus.Since {
				t.Errorf("Expected since %q for unit %q in relation %q, got %q", expectedStatus.Since, unit, relation, status.Since)
			}
		}
	}

	if fakeRunner.Command != "goal-state" {
		t.Errorf("Expected command %q, got %q", "goal-state", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}
