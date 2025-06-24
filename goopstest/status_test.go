package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func SetUnitStatusActive() error {
	err := goops.SetUnitStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func SetUnitStatusBlocked() error {
	err := goops.SetUnitStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func SetUnitStatusWaiting() error {
	err := goops.SetUnitStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func SetUnitStatusMaintenance() error {
	err := goops.SetUnitStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmSetUnitStatus(t *testing.T) {
	tests := []struct {
		name                  string
		handler               func() error
		hookName              string
		expectedStatusName    goopstest.StatusName
		expectedStatusMessage string
	}{
		{
			name:               "SetUnitStatusActive",
			handler:            SetUnitStatusActive,
			hookName:           "start",
			expectedStatusName: goopstest.StatusActive,
		},
		{
			name:               "SetUnitStatusBlocked",
			handler:            SetUnitStatusBlocked,
			hookName:           "start",
			expectedStatusName: goopstest.StatusBlocked,
		},
		{
			name:               "SetUnitStatusWaiting",
			handler:            SetUnitStatusWaiting,
			hookName:           "start",
			expectedStatusName: goopstest.StatusWaiting,
		},
		{
			name:               "SetUnitStatusMaintenance",
			handler:            SetUnitStatusMaintenance,
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

func TestCharmSetUnitStatusPreset(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetUnitStatusMaintenance,
	}

	stateIn := &goopstest.State{
		UnitStatus: &goopstest.Status{
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

func SetAppStatusActive() error {
	err := goops.SetAppStatus(goops.StatusActive, "Charm is active")
	if err != nil {
		return err
	}

	return nil
}

func SetAppStatusBlocked() error {
	err := goops.SetAppStatus(goops.StatusBlocked, "This is a test message")
	if err != nil {
		return err
	}

	return nil
}

func SetAppStatusWaiting() error {
	err := goops.SetAppStatus(goops.StatusWaiting, "Waiting for something")
	if err != nil {
		return err
	}

	return nil
}

func SetAppStatusMaintenance() error {
	err := goops.SetAppStatus(goops.StatusMaintenance, "Performing maintenance")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmSetAppStatus(t *testing.T) {
	tests := []struct {
		name               string
		handler            func() error
		hookName           string
		expectedStatusName goopstest.StatusName
	}{
		{
			name:               "SetAppStatusActive",
			handler:            SetAppStatusActive,
			hookName:           "start",
			expectedStatusName: goopstest.StatusActive,
		},
		{
			name:               "SetAppStatusBlocked",
			handler:            SetAppStatusBlocked,
			hookName:           "start",
			expectedStatusName: goopstest.StatusBlocked,
		},
		{
			name:               "SetAppStatusWaiting",
			handler:            SetAppStatusWaiting,
			hookName:           "start",
			expectedStatusName: goopstest.StatusWaiting,
		},
		{
			name:               "SetAppStatusMaintenance",
			handler:            SetAppStatusMaintenance,
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
		Charm: SetAppStatusMaintenance,
	}

	stateIn := &goopstest.State{
		AppStatus: &goopstest.Status{
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

func GetUnitStatus() error {
	status, err := goops.GetUnitStatus()
	if err != nil {
		return err
	}

	if status.Name != goops.StatusActive {
		return fmt.Errorf("expected active status, got %q", status.Name)
	}

	if status.Message != "My expected message" {
		return fmt.Errorf("unexpected message: got %q, want %q", status.Message, "My expected message")
	}

	return nil
}

func TestGetUnitStatus(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetUnitStatus,
	}

	stateIn := &goopstest.State{
		UnitStatus: &goopstest.Status{
			Name:    goopstest.StatusActive,
			Message: "My expected message",
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus.Name != goopstest.StatusActive {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, goopstest.StatusActive)
	}

	if stateOut.UnitStatus.Message != "My expected message" {
		t.Errorf("got UnitStatus.Message=%q, want %q", stateOut.UnitStatus.Message, "My expected message")
	}
}

func GetUnitStatusUnknown() error {
	status, err := goops.GetUnitStatus()
	if err != nil {
		return err
	}

	if status.Name != goops.StatusUnknown {
		return fmt.Errorf("expected unknown status, got %q", status.Name)
	}

	return nil
}

func TestGetUnitStatusNotSet(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetUnitStatusUnknown,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus.Name != goopstest.StatusUnknown {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, goopstest.StatusUnknown)
	}

	if stateOut.UnitStatus.Message != "" {
		t.Errorf("got UnitStatus.Message=%q, want empty string", stateOut.UnitStatus.Message)
	}
}

func GetAppStatus() error {
	status, err := goops.GetAppStatus()
	if err != nil {
		return err
	}

	if status.Name != goops.StatusActive {
		return fmt.Errorf("expected active status, got %q", status.Name)
	}

	if status.Message != "My expected message" {
		return fmt.Errorf("unexpected message: got %q, want %q", status.Message, "My expected message")
	}

	return nil
}

func TestGetAppStatusLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetAppStatus,
	}

	stateIn := &goopstest.State{
		Leader: true,
		AppStatus: &goopstest.Status{
			Name:    goopstest.StatusActive,
			Message: "My expected message",
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got %v", ctx.CharmErr)
	}

	if stateOut.AppStatus.Name != goopstest.StatusActive {
		t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, goopstest.StatusActive)
	}

	if stateOut.AppStatus.Message != "My expected message" {
		t.Errorf("got AppStatus.Message=%q, want %q", stateOut.AppStatus.Message, "My expected message")
	}
}

func TestGetAppStatusNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetAppStatus,
	}

	stateIn := &goopstest.State{
		Leader: false,
		AppStatus: &goopstest.Status{
			Name:    goopstest.StatusActive,
			Message: "My expected message",
		},
	}

	stateOut, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected CharmErr to be set, got nil")
	}

	expectedErr := "failed to get application status: command status-get failed: ERROR finding application status: this unit is not the leader"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}

	if stateOut.AppStatus.Name != goopstest.StatusActive {
		t.Errorf("got AppStatus=%q, want %q", stateOut.AppStatus, goopstest.StatusActive)
	}

	if stateOut.AppStatus.Message != "My expected message" {
		t.Errorf("got AppStatus.Message=%q, want %q", stateOut.AppStatus.Message, "My expected message")
	}
}
