package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ActiveIfExpectedConfig() error {
	value, err := goops.GetConfigString("blabla")
	if err != nil {
		return err
	}

	if value == "expected" {
		_ = goops.SetUnitStatus(goops.StatusActive, "Config is set to expected value")
	} else {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Config is not set to expected value")
	}

	return nil
}

func ActiveInexistantConfig() error {
	_, err := goops.GetConfigString("doesntexist")
	if err != nil {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Config key does not exist")
	}

	return nil
}

func TestCharmConfig(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		key      string
		value    string
		want     string
	}{
		{
			name:     "ActiveIfExpectedConfig",
			handler:  ActiveIfExpectedConfig,
			hookName: "start",
			key:      "blabla",
			value:    "expected",
			want:     string(goops.StatusActive),
		},
		{
			name:     "BlockedIfNotExpectedConfig",
			handler:  ActiveIfExpectedConfig,
			hookName: "start",
			key:      "blabla",
			value:    "not-expected",
			want:     string(goops.StatusBlocked),
		},
		{
			name:     "ActiveInexistantConfig",
			handler:  ActiveInexistantConfig,
			hookName: "start",
			key:      "blabla",
			value:    "whatever",
			want:     string(goops.StatusBlocked),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			config := map[string]string{
				tc.key: tc.value,
			}

			stateIn := &goopstest.State{
				Config: config,
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
