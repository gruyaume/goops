package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GoalState() error {
	goalState, err := goops.GetGoalState()
	if err != nil {
		return err
	}

	if goalState.Units["example/0"].Status != goops.StatusBlocked {
		return fmt.Errorf("expected unit status 'blocked', got '%s'", goalState.Units["example/0"].Status)
	}

	if goalState.Relations["certificates"] == nil {
		return fmt.Errorf("expected relation 'certificates' to be present")
	}

	if goalState.Relations["certificates"]["tls-certificates-requirer"].Status != goops.StatusActive {
		return fmt.Errorf("expected app status 'active' for relation 'certificates:tls-certificates-requirer', got '%s'", goalState.Relations["certificates"]["tls-certificates-requirer"].Status)
	}

	if goalState.Relations["certificates"]["tls-certificates-requirer/0"].Status != goops.StatusActive {
		return fmt.Errorf("expected app status 'active' for relation 'certificates:tls-certificates-requirer/0', got '%s'", goalState.Relations["certificates"]["tls-certificates-requirer/0"].Status)
	}

	return nil
}

func TestGoalState(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GoalState,
		AppName: "example",
		UnitID:  "example/0",
	}

	stateIn := goopstest.State{
		UnitStatus: goopstest.Status{
			Name:    goopstest.StatusBlocked,
			Message: "Unit is blocked",
		},
		Relations: []goopstest.Relation{
			{
				Endpoint:      "certificates",
				RemoteAppName: "tls-certificates-requirer",
				RemoteUnitsData: map[goopstest.UnitID]goopstest.DataBag{
					"tls-certificates-requirer/0": {
						"ca": "example-ca-cert",
					},
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}
}
