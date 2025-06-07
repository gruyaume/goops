package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func MaintenanceStatusOnAction() error {
	env := goops.ReadEnv()

	if env.ActionName == "run-action" {
		err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
		if err != nil {
			return err
		}
	} else {
		err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
		if err != nil {
			return err
		}
	}

	return nil
}

func TestCharmActionName(t *testing.T) {
	tests := []struct {
		name       string
		handler    func() error
		actionName string
		want       string
	}{
		{
			name:       "MaintenanceStatusOnAction",
			handler:    MaintenanceStatusOnAction,
			actionName: "run-action",
			want:       string(goops.StatusMaintenance),
		},
		{
			name:       "ActiveStatusOnOtherActions",
			handler:    MaintenanceStatusOnAction,
			actionName: "something-else",
			want:       string(goops.StatusActive),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			stateIn := &goopstest.State{}

			stateOut, err := ctx.RunAction(tc.actionName, stateIn)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if stateOut.UnitStatus != tc.want {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.want)
			}
		})
	}
}
