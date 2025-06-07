package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func MaintenanceStatusOnStart() error {
	env := goops.ReadEnv()
	if env.HookName == "start" {
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

func ActiveStatusOnInstall() error {
	env := goops.ReadEnv()
	if env.HookName == "install" {
		err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
		if err != nil {
			return err
		}
	} else {
		err := goops.SetUnitStatus(goops.StatusBlocked, "Charm is blocked")
		if err != nil {
			return err
		}
	}

	return nil
}

func TestCharmHookName(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		want     string
	}{
		{
			name:     "MaintenanceStatusOnStart",
			handler:  MaintenanceStatusOnStart,
			hookName: "start",
			want:     string(goops.StatusMaintenance),
		},
		{
			name:     "ActiveStatusOnInstall",
			handler:  ActiveStatusOnInstall,
			hookName: "install",
			want:     string(goops.StatusActive),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			stateIn := &goopstest.State{}

			stateOut, err := ctx.Run(tc.hookName, stateIn)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if stateOut.UnitStatus != tc.want {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.want)
			}
		})
	}
}
