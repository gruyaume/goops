package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ActiveStatus() error {
	err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func BlockedStatus() error {
	err := goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func WaitingStatus() error {
	err := goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func MaintenanceStatus() error {
	err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmStatus(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		want     string
	}{
		{
			name:     "ActiveStatus",
			handler:  ActiveStatus,
			hookName: "start",
			want:     string(goops.StatusActive),
		},
		{
			name:     "BlockedStatus",
			handler:  BlockedStatus,
			hookName: "start",
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "WaitingStatus",
			handler:  WaitingStatus,
			hookName: "start",
			want:     string(goops.StatusWaiting),
		},
		{
			name:     "MaintenanceStatus",
			handler:  MaintenanceStatus,
			hookName: "start",
			want:     string(goops.StatusMaintenance),
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

func TestCharmStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: MaintenanceStatus,
	}

	stateIn := &goopstest.State{
		UnitStatus: string(goops.StatusActive),
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus != string(goops.StatusMaintenance) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusMaintenance))
	}
}
