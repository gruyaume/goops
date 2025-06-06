package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ConfigureActive() error {
	_ = goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	return nil
}

func ConfigureBlocked() error {
	_ = goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	return nil
}

func ConfigureWaiting() error {
	_ = goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	return nil
}

func ConfigureMaintenance() error {
	_ = goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	return nil
}

func ConfigureMaintenanceOnInstall() error {
	env := goops.ReadEnv()
	if env.HookName == "start" {
		_ = goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	} else {
		_ = goops.SetUnitStatus(goops.StatusActive, "Charm is active")
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
			handler:  ConfigureActive,
			hookName: "start",
			want:     string(goops.StatusActive),
		},
		{
			name:     "BlockedStatus",
			handler:  ConfigureBlocked,
			hookName: "start",
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "WaitingStatus",
			handler:  ConfigureWaiting,
			hookName: "start",
			want:     string(goops.StatusWaiting),
		},
		{
			name:     "MaintenanceStatus",
			handler:  ConfigureMaintenance,
			hookName: "start",
			want:     string(goops.StatusMaintenance),
		},
		{
			name:     "MaintenanceStatusOnInstall",
			handler:  ConfigureMaintenanceOnInstall,
			hookName: "start",
			want:     string(goops.StatusMaintenance),
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			stateIn := goopstest.State{}

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
