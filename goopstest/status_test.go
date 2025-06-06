package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ConfigureActive() error {
	err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func ConfigureBlocked() error {
	err := goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func ConfigureWaiting() error {
	err := goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func ConfigureMaintenance() error {
	err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func ConfigureMaintenanceOnStart() error {
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

func ConfigureMaintenanceIfLeader() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return err
	}

	if isLeader {
		err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance on leader")
		if err != nil {
			return err
		}
	} else {
		err := goops.SetUnitStatus(goops.StatusActive, "Charm is active on non-leader")
		if err != nil {
			return err
		}
	}

	return nil
}

func ConfigureMaintenanceIfNotLeader() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return err
	}

	if !isLeader {
		err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance on leader")
		if err != nil {
			return err
		}
	} else {
		err := goops.SetUnitStatus(goops.StatusActive, "Charm is active on non-leader")
		if err != nil {
			return err
		}
	}

	return nil
}

func TestCharmStatus(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		leader   bool
		want     string
	}{
		{
			name:     "ActiveStatus",
			handler:  ConfigureActive,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusActive),
		},
		{
			name:     "BlockedStatus",
			handler:  ConfigureBlocked,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "WaitingStatus",
			handler:  ConfigureWaiting,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusWaiting),
		},
		{
			name:     "MaintenanceStatus",
			handler:  ConfigureMaintenance,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusMaintenance),
		},
		{
			name:     "MaintenanceStatusOnStart",
			handler:  ConfigureMaintenanceOnStart,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusMaintenance),
		},
		{
			name:     "MaintenanceStatusIfLeader",
			handler:  ConfigureMaintenanceIfLeader,
			hookName: "start",
			leader:   true,
			want:     string(goops.StatusMaintenance),
		},
		{
			name:     "MaintenanceStatusIfNotLeader",
			handler:  ConfigureMaintenanceIfNotLeader,
			hookName: "start",
			leader:   false,
			want:     string(goops.StatusMaintenance),
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			stateIn := &goopstest.State{
				Leader: tc.leader,
			}

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
