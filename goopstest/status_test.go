package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func UnitActiveStatus() error {
	err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func UnitBlockedStatus() error {
	err := goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func UnitWaitingStatus() error {
	err := goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func UnitMaintenanceStatus() error {
	err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmUnitStatus(t *testing.T) {
	tests := []struct {
		name                  string
		handler               func() error
		hookName              string
		expectedStatusName    goopstest.StatusName
		expectedStatusMessage string
	}{
		{
			name:               "UnitActiveStatus",
			handler:            UnitActiveStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusActive,
		},
		{
			name:               "UnitBlockedStatus",
			handler:            UnitBlockedStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusBlocked,
		},
		{
			name:               "UnitWaitingStatus",
			handler:            UnitWaitingStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusWaiting,
		},
		{
			name:               "UnitMaintenanceStatus",
			handler:            UnitMaintenanceStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusMaintenance,
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

			if stateOut.UnitStatus.Name != tc.expectedStatusName {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.expectedStatusName)
			}
		})
	}
}

func TestCharmUnitStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: UnitMaintenanceStatus,
	}

	stateIn := &goopstest.State{
		UnitStatus: goopstest.Status{
			Name: goopstest.StatusActive,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus.Name != goopstest.StatusMaintenance {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, goopstest.StatusMaintenance)
	}

	if stateOut.UnitStatus.Message != "Performing maintenance" {
		t.Errorf("got UnitStatus.Message=%q, want %q", stateOut.UnitStatus.Message, "Performing maintenance")
	}
}

func AppActiveStatus() error {
	err := goops.SetAppStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func AppBlockedStatus() error {
	err := goops.SetAppStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func AppWaitingStatus() error {
	err := goops.SetAppStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func AppMaintenanceStatus() error {
	err := goops.SetAppStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmAppStatus(t *testing.T) {
	tests := []struct {
		name               string
		handler            func() error
		hookName           string
		expectedStatusName goopstest.StatusName
	}{
		{
			name:               "AppActiveStatus",
			handler:            AppActiveStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusActive,
		},
		{
			name:               "AppBlockedStatus",
			handler:            AppBlockedStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusBlocked,
		},
		{
			name:               "AppWaitingStatus",
			handler:            AppWaitingStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusWaiting,
		},
		{
			name:               "AppMaintenanceStatus",
			handler:            AppMaintenanceStatus,
			hookName:           "start",
			expectedStatusName: goopstest.StatusMaintenance,
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

			if stateOut.AppStatus.Name != tc.expectedStatusName {
				t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, tc.expectedStatusName)
			}
		})
	}
}

func TestCharmAppStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: AppMaintenanceStatus,
	}

	stateIn := &goopstest.State{
		AppStatus: goopstest.Status{
			Name: goopstest.StatusActive,
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.AppStatus.Name != goopstest.StatusMaintenance {
		t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, goopstest.StatusMaintenance)
	}

	if stateOut.AppStatus.Message != "Performing maintenance" {
		t.Errorf("got AppStatus.Message=%q, want %q", stateOut.AppStatus.Message, "Performing maintenance")
	}
}
