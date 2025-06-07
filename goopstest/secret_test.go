package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ConfigureGetSecret() error {
	secretLabel := "whatever-label"

	secret, err := goops.GetSecretByLabel(secretLabel, false, true)
	if err != nil {
		_ = goops.SetUnitStatus(goops.StatusMaintenance, "Secret does not exist")
		return nil
	}

	secretValue, ok := secret["secret-key"]
	if !ok {
		return goops.FailActionf("secret key not found in secret with label %s", secretLabel)
	}

	if secretValue != "expected-value" {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Secret is not set to expected value")
		return nil
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Secret is set to expected value")

	return nil
}

func TestCharmSecret(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		key      string
		value    string
		want     string
	}{
		{
			name:     "ConfigureGetSecret",
			handler:  ConfigureGetSecret,
			hookName: "start",
			key:      "secret-key",
			value:    "expected-value",
			want:     string(goops.StatusActive),
		},
		{
			name:     "ConfigureGetSecretUnexpectedValue",
			handler:  ConfigureGetSecret,
			hookName: "start",
			key:      "secret-key",
			value:    "unexpected-value",
			want:     string(goops.StatusBlocked),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := goopstest.Context{
				Charm: tc.handler,
			}

			mySecret := goopstest.Secret{
				Label: "whatever-label",
				Content: map[string]string{
					tc.key: tc.value,
				},
			}

			stateIn := &goopstest.State{
				Secrets: []goopstest.Secret{
					mySecret,
				},
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

func TestCharmUnexistingSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ConfigureGetSecret,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus != string(goops.StatusMaintenance) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusMaintenance))
	}
}
