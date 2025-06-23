package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

type MyConfig struct {
	WhateverKey string `json:"whatever_key"`
}

func ActiveIfExpectedConfig() error {
	myConfig := MyConfig{}

	err := goops.GetConfig(&myConfig)
	if err != nil {
		return err
	}

	if myConfig.WhateverKey == "expected" {
		_ = goops.SetUnitStatus(goops.StatusActive, "Config is set to expected value")
	} else {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Config is not set to expected value")
	}

	return nil
}

type MyBadConfig struct {
	WhateverKey string `json:"whatever_key"`
}

func ActiveInexistantConfig() error {
	myBadConfig := MyBadConfig{}

	err := goops.GetConfig(&myBadConfig)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if myBadConfig.WhateverKey == "" {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Config is not set")
	}

	return nil
}

func TestCharmConfig(t *testing.T) {
	tests := []struct {
		name               string
		handler            func() error
		hookName           string
		key                string
		value              string
		expectedStatusName goopstest.StatusName
	}{
		{
			name:               "ActiveIfExpectedConfig",
			handler:            ActiveIfExpectedConfig,
			hookName:           "start",
			key:                "whatever_key",
			value:              "expected",
			expectedStatusName: goopstest.StatusActive,
		},
		{
			name:               "BlockedIfNotExpectedConfig",
			handler:            ActiveIfExpectedConfig,
			hookName:           "start",
			key:                "whatever_key",
			value:              "not-expected",
			expectedStatusName: goopstest.StatusBlocked,
		},
		{
			name:               "ActiveInexistantConfig",
			handler:            ActiveInexistantConfig,
			hookName:           "start",
			key:                "a-different-key",
			value:              "whatever",
			expectedStatusName: goopstest.StatusBlocked,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			config := map[string]any{
				tc.key: tc.value,
			}

			stateIn := &goopstest.State{
				Config: config,
			}

			stateOut, err := ctx.Run(tc.hookName, stateIn)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if ctx.CharmErr != nil {
				t.Errorf("expected no error, got %v", ctx.CharmErr)
			}

			if stateOut.UnitStatus.Name != tc.expectedStatusName {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.expectedStatusName)
			}
		})
	}
}

func TestActiveIfExpectedConfigInActionHook(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ActiveIfExpectedConfig,
	}

	config := map[string]any{
		"whatever_key": "expected",
	}

	stateIn := &goopstest.State{
		Config: config,
	}

	stateOut, err := ctx.RunAction("run-action", stateIn, nil)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus.Name != goopstest.StatusActive {
		t.Errorf("Expected UnitStatus %q, got %q", goopstest.StatusActive, stateOut.UnitStatus)
	}
}
