package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetSecretByLabel() error {
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

func TestCharmGetSecretByLabel(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		key      string
		value    string
		want     string
	}{
		{
			name:     "GetSecretByLabel",
			handler:  GetSecretByLabel,
			hookName: "start",
			key:      "secret-key",
			value:    "expected-value",
			want:     string(goops.StatusActive),
		},
		{
			name:     "GetSecretByLabelUnexpectedValue",
			handler:  GetSecretByLabel,
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

			mySecret := &goopstest.Secret{
				Label: "whatever-label",
				Content: map[string]string{
					tc.key: tc.value,
				},
			}

			stateIn := &goopstest.State{
				Secrets: []*goopstest.Secret{
					mySecret,
				},
				Leader: true,
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

func TestCharmGetUnexistingSecretByLabel(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretByLabel,
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

func GetSecretByID() error {
	secretID := "12345"

	secret, err := goops.GetSecretByID(secretID, false, true)
	if err != nil {
		return err
	}

	secretValue, ok := secret["secret-key"]
	if !ok {
		return goops.FailActionf("secret key not found in secret with ID %s", secretID)
	}

	if secretValue != "expected-value" {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Secret is not set to expected value")
		return nil
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Secret is set to expected value")

	return nil
}

func TestCharmGetSecretByID(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() error
		hookName string
		key      string
		value    string
		want     string
	}{
		{
			name:     "GetSecretByID",
			handler:  GetSecretByID,
			hookName: "start",
			key:      "secret-key",
			value:    "expected-value",
			want:     string(goops.StatusActive),
		},
		{
			name:     "GetSecretByIDUnexpectedValue",
			handler:  GetSecretByID,
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

			mySecret := &goopstest.Secret{
				ID: "12345",
				Content: map[string]string{
					tc.key: tc.value,
				},
			}

			stateIn := &goopstest.State{
				Secrets: []*goopstest.Secret{
					mySecret,
				},
				Leader: true,
			}

			stateOut, err := ctx.Run(tc.hookName, stateIn)
			if err != nil {
				t.Fatalf("Run returned an error: %v", err)
			}

			if ctx.CharmErr != nil {
				t.Fatalf("Run returned an error: %v", ctx.CharmErr)
			}

			if stateOut.UnitStatus != tc.want {
				t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, tc.want)
			}

			for _, secret := range stateOut.Secrets {
				if secret.ID == mySecret.ID {
					if secret.Content[tc.key] != tc.value {
						t.Errorf("got Secret[%s]=%s, want %s", tc.key, secret.Content[tc.key], tc.value)
					}
				}
			}
		})
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
		return err
	}

	return nil
}

func TestCharmAddSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: AddSecret,
	}

	stateIn := &goopstest.State{
		Leader: true,
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Secrets) != 1 {
		t.Fatalf("got %d secrets, want 1", len(stateOut.Secrets))
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

func TestCharmAddSecretNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: AddSecret,
	}

	stateIn := &goopstest.State{
		Leader: false,
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected an error when not leader, got nil")
	}

	if ctx.CharmErr.Error() != "failed to add secret: command secret-add failed: ERROR this unit is not the leader" {
		t.Errorf("got CharmErr=%q, want 'failed to add secret: command secret-add failed: ERROR this unit is not the leader'", ctx.CharmErr.Error())
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
		Leader: true,
		Secrets: []*goopstest.Secret{
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

func TestCharmRemoveSecretNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RemoveSecret,
	}

	stateIn := &goopstest.State{
		Leader: false,
		Secrets: []*goopstest.Secret{
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

	if len(stateOut.Secrets) != 1 {
		t.Fatalf("got %d secrets, want 1", len(stateOut.Secrets))
	}

	if stateOut.Secrets[0].ID != "123" {
		t.Errorf("got Secret.ID=%s, want %s", stateOut.Secrets[0].ID, "123")
	}
}

func TestCharmRemoveUnexistingSecret(t *testing.T) {
	ctx := goopstest.Context{
		Charm: RemoveSecret,
	}

	stateIn := &goopstest.State{
		Secrets: []*goopstest.Secret{
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

func GetSecretInfoByID() error {
	secretID := "12345"

	secretInfo, err := goops.GetSecretInfoByID(secretID)
	if err != nil {
		return err
	}

	if len(secretInfo) == 0 {
		return fmt.Errorf("no secret info found for ID: %s", secretID)
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Secret info retrieved successfully")

	return nil
}

func TestCharmGetSecretInfoByID(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretInfoByID,
	}

	stateIn := &goopstest.State{
		Secrets: []*goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
		},
		Leader: true,
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("Run returned an error: %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusActive))
	}
}

func TestCharmGetSecretInfoByIDNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretInfoByID,
	}

	stateIn := &goopstest.State{
		Leader: false,
		Secrets: []*goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected an error when not leader, got nil")
	}

	if ctx.CharmErr.Error() != `failed to get secret info: ERROR secret "12345" not found` {
		t.Errorf("got CharmErr=%q, want 'failed to get secret info: ERROR secret \"12345\" not found'", ctx.CharmErr.Error())
	}
}

func GetSecretInfoByLabel() error {
	secretLabel := "whatever-label"

	secretInfo, err := goops.GetSecretInfoByLabel(secretLabel)
	if err != nil {
		return fmt.Errorf("could not get secret info: %w", err)
	}

	if len(secretInfo) == 0 {
		return fmt.Errorf("no secret info found for label: %s", secretLabel)
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Secret info retrieved successfully")

	return nil
}

func TestCharmGetSecretInfoByLabel(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretInfoByLabel,
	}

	stateIn := &goopstest.State{
		Leader: true,
		Secrets: []*goopstest.Secret{
			{
				Label: "whatever-label",
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

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusActive))
	}
}

func GetSecretIDs() error {
	secretIDs, err := goops.GetSecretIDs()
	if err != nil {
		return fmt.Errorf("could not get secret IDs: %w", err)
	}

	if len(secretIDs) != 2 {
		return fmt.Errorf("expected 2 secret IDs, got %d", len(secretIDs))
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Secret IDs retrieved successfully")

	return nil
}

func TestCharmGetSecretIDs(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretIDs,
	}

	stateIn := &goopstest.State{
		Leader: true,
		Secrets: []*goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
			{
				ID: "67890",
				Content: map[string]string{
					"private-key":    "another-keycontent",
					"ca-certificate": "another-certcontent",
				},
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("Run returned an error: %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusActive))
	}

	if len(stateOut.Secrets) != 2 {
		t.Errorf("got %d secrets, want 2", len(stateOut.Secrets))
	}
}

func TestCharmGetSecretIDsNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSecretIDs,
	}

	stateIn := &goopstest.State{
		Leader: false,
		Secrets: []*goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
			{
				ID: "67890",
				Content: map[string]string{
					"private-key":    "another-keycontent",
					"ca-certificate": "another-certcontent",
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected an error when not leader, got nil")
	}

	if ctx.CharmErr.Error() != "expected 2 secret IDs, got 0" {
		t.Errorf("got CharmErr=%q, want 'expected 2 secret IDs, got 0'", ctx.CharmErr.Error())
	}
}

func GrantSecretToRelation() error {
	secretID := "12345"
	relation := "db:2"

	err := goops.GrantSecretToRelation(secretID, relation)
	if err != nil {
		return err
	}

	_ = goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Secret %s granted to relation %s", secretID, relation))

	return nil
}

func TestCharmGrantSecretToRelation(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GrantSecretToRelation,
	}

	stateIn := &goopstest.State{
		Leader: true,
		Secrets: []*goopstest.Secret{
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

	if ctx.CharmErr != nil {
		t.Fatalf("Run returned an error: %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusActive))
	}
}

func GrantSecretToUnit() error {
	secretID := "12345"
	relation := "db:0"
	unitName := "db/1"

	err := goops.GrantSecretToUnit(secretID, relation, unitName)
	if err != nil {
		return err
	}

	_ = goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Secret %s granted to unit %s", secretID, unitName))

	return nil
}

func TestCharmGrantSecretToUnit(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GrantSecretToUnit,
	}

	stateIn := &goopstest.State{
		Leader: true,
		Secrets: []*goopstest.Secret{
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

	if ctx.CharmErr != nil {
		t.Fatalf("Run returned an error: %v", ctx.CharmErr)
	}

	if stateOut.UnitStatus != string(goops.StatusActive) {
		t.Errorf("got UnitStatus=%q, want %q", stateOut.UnitStatus, string(goops.StatusActive))
	}
}

func TestCharmGrantSecretNonLeader(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GrantSecretToRelation,
	}

	stateIn := &goopstest.State{
		Leader: false,
		Secrets: []*goopstest.Secret{
			{
				ID: "12345",
				Content: map[string]string{
					"private-key":    "keycontent",
					"ca-certificate": "certcontent",
				},
			},
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected an error when not leader, got nil")
	}

	if ctx.CharmErr.Error() != "failed to grant secret: ERROR secret \"12345\" not found" {
		t.Errorf("got CharmErr=%q, want 'failed to grant secret: ERROR secret \"12345\" not found'", ctx.CharmErr.Error())
	}
}
