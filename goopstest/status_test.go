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
		name     string
		handler  func() error
		hookName string
		want     string
	}{
		{
			name:     "UnitActiveStatus",
			handler:  UnitActiveStatus,
			hookName: "start",
			want:     string(goops.StatusActive),
		},
		{
			name:     "UnitBlockedStatus",
			handler:  UnitBlockedStatus,
			hookName: "start",
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "UnitWaitingStatus",
			handler:  UnitWaitingStatus,
			hookName: "start",
			want:     string(goops.StatusWaiting),
		},
		{
			name:     "UnitMaintenanceStatus",
			handler:  UnitMaintenanceStatus,
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

func TestCharmUnitStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: UnitMaintenanceStatus,
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
		name     string
		handler  func() error
		hookName string
		want     string
	}{
		{
			name:     "AppActiveStatus",
			handler:  AppActiveStatus,
			hookName: "start",
			want:     string(goops.StatusActive),
		},
		{
			name:     "AppBlockedStatus",
			handler:  AppBlockedStatus,
			hookName: "start",
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "AppWaitingStatus",
			handler:  AppWaitingStatus,
			hookName: "start",
			want:     string(goops.StatusWaiting),
		},
		{
			name:     "AppMaintenanceStatus",
			handler:  AppMaintenanceStatus,
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

			if stateOut.AppStatus != tc.want {
				t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, tc.want)
			}
		})
	}
}

func TestCharmAppStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: AppMaintenanceStatus,
	}

	stateIn := &goopstest.State{
		AppStatus: string(goops.StatusActive),
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.AppStatus != string(goops.StatusMaintenance) {
		t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, string(goops.StatusMaintenance))
	}
}
