package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func MaintenanceStatusIfLeader() error {
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

func MaintenanceStatusIfNotLeader() error {
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
		name               string
		handler            func() error
		hookName           string
		leader             bool
		expectedStatusName goopstest.StatusName
	}{
		{
			name:               "MaintenanceStatusIfLeader",
			handler:            MaintenanceStatusIfLeader,
			hookName:           "start",
			leader:             true,
			expectedStatusName: goopstest.StatusMaintenance,
		},
		{
			name:               "MaintenanceStatusIfNotLeader",
			handler:            MaintenanceStatusIfNotLeader,
			hookName:           "start",
			leader:             false,
			expectedStatusName: goopstest.StatusMaintenance,
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

			if stateOut.UnitStatus.Name != tc.expectedStatusName {
				t.Errorf("got Status=%v, want %v", stateOut.UnitStatus, tc.expectedStatusName)
			}

			if stateOut.Leader != tc.leader {
				t.Errorf("got Leader=%v, want %v", stateOut.Leader, tc.leader)
			}
		})
	}
}
