package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetSecret() error {
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

func TestCharmGetSecret(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		key      string
		value    string
		want     string
	}{
		{
			name:     "GetSecret",
			handler:  GetSecret,
			hookName: "start",
			key:      "secret-key",
			value:    "expected-value",
			want:     string(goops.StatusActive),
		},
		{
			name:     "GetSecretUnexpectedValue",
			handler:  GetSecret,
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

			for _, secret := range stateOut.Secrets {
				if secret.Label == mySecret.Label {
					if secret.Content[tc.key] != tc.value {
						t.Errorf("got Secret[%s]=%s, want %s", tc.key, secret.Content[tc.key], tc.value)
					}
				}
			}
		})
	}
}

func TestCharmGetUnexistingSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecret,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.UnitStatus != string(goops.StatusMaintenance) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusMaintenance))
	}

	if len(stateOut.Secrets) != 0 {
		t.Errorf("got %d secrets, want 0", len(stateOut.Secrets))
	}
}

func AddSecret() error {
	secretLabel := "whatever-label"

	caKeyPEM := `keycontent`
	caCertPEM := `certcontent`

	secretContent := map[string]string{
		"private-key":    caKeyPEM,
		"ca-certificate": caCertPEM,
	}

	_, err := goops.AddSecret(&goops.AddSecretOptions{
		Label:   secretLabel,
		Content: secretContent,
	})
	if err != nil {
		return fmt.Errorf("could not add secret: %w", err)
	}

	return nil
}

func TestCharmAddSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: AddSecret,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Secrets) != 1 {
		t.Errorf("got %d secrets, want 1", len(stateOut.Secrets))
	}

	mySecret := stateOut.Secrets[0]
	if mySecret.Label != "whatever-label" {
		t.Errorf("got Secret.Label=%s, want %s", mySecret.Label, "whatever-label")
	}

	if mySecret.Content["private-key"] != "keycontent" {
		t.Errorf("got Secret[private-key]=%s, want %s", mySecret.Content["private-key"], "keycontent")
	}

	if mySecret.Content["ca-certificate"] != "certcontent" {
		t.Errorf("got Secret[ca-certificate]=%s, want %s", mySecret.Content["ca-certificate"], "certcontent")
	}
}

func RemoveSecret() error {
	err := goops.RemoveSecret("123")
	if err != nil {
		return err
	}

	return nil
}

func TestCharmRemoveSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RemoveSecret,
	}

	stateIn := &goopstest.State{
		Secrets: []goopstest.Secret{
			{
				ID: "123",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Secrets) != 0 {
		t.Errorf("got %d secrets, want 0", len(stateOut.Secrets))
	}
}

func TestCharmRemoveUnexistingSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RemoveSecret,
	}

	stateIn := &goopstest.State{
		Secrets: []goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Secrets) != 1 {
		t.Errorf("got %d secrets, want 1", len(stateOut.Secrets))
	}

	if stateOut.Secrets[0].ID != "12345" {
		t.Errorf("got Secret.ID=%s, want %s", stateOut.Secrets[0].ID, "12345")
	}
}
