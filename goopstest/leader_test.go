package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

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

func TestCharmLeader(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		leader   bool
		want     string
	}{
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
				t.Errorf("got Status=%v, want %v", stateOut.UnitStatus, tc.want)
			}

			if stateOut.Leader != tc.leader {
				t.Errorf("got Leader=%v, want %v", stateOut.Leader, tc.leader)
			}
		})
	}
}
