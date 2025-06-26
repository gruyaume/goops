package goopstest_test

import (
	"fmt"
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
		name               string
		handler            func() error
		actionName         string
		expectedStatusName goopstest.StatusName
	}{
		{
			name:               "MaintenanceStatusOnAction",
			handler:            MaintenanceStatusOnAction,
			actionName:         "run-action",
			expectedStatusName: goopstest.StatusMaintenance,
		},
		{
			name:               "ActiveStatusOnOtherActions",
			handler:            MaintenanceStatusOnAction,
			actionName:         "something-else",
			expectedStatusName: goopstest.StatusActive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.NewContext(tc.handler)

			stateIn := goopstest.State{}

			stateOut, err := ctx.RunAction(tc.actionName, stateIn, nil)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if stateOut.UnitStatus.Name != tc.expectedStatusName {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.expectedStatusName)
			}
		})
	}
}

func ActionResults1() error {
	results := map[string]string{
		"key": "value",
	}

	err := goops.SetActionResults(results)
	if err != nil {
		return err
	}

	return nil
}

func TestCharmActionResults1(t *testing.T) {
	ctx := goopstest.NewContext(ActionResults1)

	stateIn := goopstest.State{}

	_, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.ActionResults["key"] != "value" {
		t.Errorf("got ActionResults[key]=%s, want value", ctx.ActionResults["key"])
	}
}

func ActionResults3() error {
	results := map[string]string{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
	}

	err := goops.SetActionResults(results)
	if err != nil {
		return err
	}

	return nil
}

func TestCharmActionResults3(t *testing.T) {
	ctx := goopstest.NewContext(ActionResults3)

	stateIn := goopstest.State{}

	_, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.ActionResults["key0"] != "value0" {
		t.Errorf("got ActionResults[key0]=%s, want value0", ctx.ActionResults["key0"])
	}

	if ctx.ActionResults["key1"] != "value1" {
		t.Errorf("got ActionResults[key1]=%s, want value1", ctx.ActionResults["key1"])
	}

	if ctx.ActionResults["key2"] != "value2" {
		t.Errorf("got ActionResults[key2]=%s, want value2", ctx.ActionResults["key2"])
	}
}

func ActionFailed() error {
	err := goops.FailActionf("Action failed with error: %s", "some error")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmActionFailed(t *testing.T) {
	ctx := goopstest.NewContext(ActionFailed)

	stateIn := goopstest.State{}

	_, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.ActionError.Error() != "Action failed with error: some error" {
		t.Errorf("got ActionError=%q, want %q", ctx.ActionError.Error(), "Action failed with error: some error")
	}
}

type ExampleActionParams struct {
	Email     string `json:"email"`
	AcceptTOS bool   `json:"accept-tos"`
}

func GetActionParamsAndSetResults() error {
	actionParams := ExampleActionParams{}

	err := goops.GetActionParams(&actionParams)
	if err != nil {
		_ = goops.FailActionf("couldn't get action parameters: %v", err)
		return nil
	}

	if actionParams.Email != "expected-value" {
		_ = goops.FailActionf("Action parameter 'whatever-key' not set")
		return nil
	}

	if !actionParams.AcceptTOS {
		_ = goops.FailActionf("You must accept the terms of service to run this action")
		return nil
	}

	err = goops.SetActionResults(map[string]string{
		"success": "true",
	})
	if err != nil {
		return err
	}

	return nil
}

func TestCharmActionParameters(t *testing.T) {
	ctx := goopstest.NewContext(GetActionParamsAndSetResults)

	stateIn := goopstest.State{}

	_, err := ctx.RunAction("run-action", stateIn, map[string]any{
		"email":      "expected-value",
		"accept-tos": true,
	})
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.ActionResults["success"] != "true" {
		t.Errorf("got ActionResults[success]=%s, want true", ctx.ActionResults["success"])
	}
}

func TestCharmActionParameterNotSet(t *testing.T) {
	ctx := goopstest.NewContext(GetActionParamsAndSetResults)

	stateIn := goopstest.State{}

	_, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.ActionError == nil {
		t.Fatal("Expected ActionError to be set, got nil")
	}

	if ctx.ActionError.Error() != "Action parameter 'whatever-key' not set" {
		t.Errorf("got ActionError=%q, want 'Action parameter 'whatever-key' not set'", ctx.ActionError.Error())
	}

	if ctx.ActionResults != nil {
		t.Errorf("got ActionResults=%v, want nil", ctx.ActionResults)
	}
}

func ActionParams() error {
	actionParams := ExampleActionParams{}

	err := goops.GetActionParams(&actionParams)
	if err != nil {
		return fmt.Errorf("couldn't get action parameters: %w", err)
	}

	return nil
}

func TestGetActionParamInNonActionHook(t *testing.T) {
	ctx := goopstest.NewContext(ActionParams)

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	if ctx.CharmErr.Error() != "couldn't get action parameters: failed to get action parameter: command action-get failed: ERROR not running an action" {
		t.Errorf("got CharmErr=%q, want 'couldn't get action parameters: failed to get action parameter: command action-get failed: ERROR not running an action'", ctx.CharmErr.Error())
	}
}

func ActionFailf() error {
	err := goops.FailActionf("whatever message")
	if err != nil {
		return err
	}

	return nil
}

func TestActionFailfInNonActionHook(t *testing.T) {
	ctx := goopstest.NewContext(ActionFailf)

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	if ctx.CharmErr.Error() != "failed to fail action: command action-fail failed: ERROR not running an action" {
		t.Errorf("got CharmErr=%q, want 'failed to fail action: command action-fail failed: ERROR not running an action'", ctx.CharmErr.Error())
	}
}

func ActionLogf() error {
	err := goops.ActionLogf("This is a log message from the action")
	if err != nil {
		return err
	}

	return nil
}

func TestActionLogfInNonActionHook(t *testing.T) {
	ctx := goopstest.NewContext(ActionLogf)

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	if ctx.CharmErr.Error() != "failed to log action message: command action-log failed: ERROR not running an action" {
		t.Errorf("got CharmErr=%q, want 'failed to log action message: command action-log failed: ERROR not running an action'", ctx.CharmErr.Error())
	}
}

func SetActionResults() error {
	results := map[string]string{
		"key": "value",
	}

	err := goops.SetActionResults(results)
	if err != nil {
		return err
	}

	return nil
}

func TestSetActionResultsInNonActionHook(t *testing.T) {
	ctx := goopstest.NewContext(SetActionResults)

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr == nil {
		t.Fatal("Expected CharmErr to be set, got nil")
	}

	if ctx.CharmErr.Error() != "failed to set action parameters: command action-set failed: ERROR not running an action" {
		t.Errorf("got CharmErr=%q, want 'failed to set action parameters: command action-set failed: ERROR not running an action'", ctx.CharmErr.Error())
	}
}
