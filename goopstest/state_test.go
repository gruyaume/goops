package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetState() error {
	value, err := goops.GetState("my-key")
	if err != nil {
		return err
	}

	if value != "my-value" {
		return fmt.Errorf("unexpected state value: got %s, want my-value", value)
	}

	return nil
}

func TestGetState(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetState,
	}

	stateIn := goopstest.State{
		StoredState: map[string]string{
			"my-key": "my-value",
		},
	}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got: %v", ctx.CharmErr)
	}
}

func TestGetStateNoKey(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetState,
	}

	stateIn := goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedErr := "no state found for key: my-key"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("got CharmErr=%q, want %q", ctx.CharmErr.Error(), expectedErr)
	}
}

func SetState() error {
	err := goops.SetState("my-key", "my-value")
	if err != nil {
		return fmt.Errorf("could not set state: %w", err)
	}

	return nil
}

func TestSetState(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetState,
	}

	stateIn := goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.StoredState["my-key"] != "my-value" {
		t.Errorf("got StoredState[my-key]=%s, want my-value", stateOut.StoredState["my-key"])
	}
}

func GetSetState() error {
	initialValue, err := goops.GetState("my-key")
	if err != nil {
		return fmt.Errorf("could not get state: %w", err)
	}

	if initialValue != "my-value" {
		return fmt.Errorf("unexpected initial state value: got %s, want my-value", initialValue)
	}

	err = goops.SetState("my-key", "my-new-value")
	if err != nil {
		return fmt.Errorf("could not set state: %w", err)
	}

	newValue, err := goops.GetState("my-key")
	if err != nil {
		return fmt.Errorf("could not get new state: %w", err)
	}

	if newValue != "my-new-value" {
		return fmt.Errorf("unexpected new state value: got %s, want my-new-value", newValue)
	}

	return nil
}

func TestGetSetState(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetSetState,
	}

	stateIn := goopstest.State{
		StoredState: map[string]string{
			"my-key": "my-value",
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.StoredState["my-key"] != "my-new-value" {
		t.Errorf("got StoredState[my-key]=%s, want my-new-value", stateOut.StoredState["my-key"])
	}
}

func DeleteState() error {
	err := goops.DeleteState("my-key")
	if err != nil {
		return err
	}

	return nil
}

func TestDeleteState(t *testing.T) {
	ctx := goopstest.Context{
		Charm: DeleteState,
	}

	stateIn := goopstest.State{
		StoredState: map[string]string{
			"my-key": "my-value",
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if _, exists := stateOut.StoredState["my-key"]; exists {
		t.Errorf("expected StoredState[my-key] to be deleted, but it still exists")
	}

	if len(stateOut.StoredState) != 0 {
		t.Errorf("expected StoredState to be empty, got %d items", len(stateOut.StoredState))
	}
}

func TestDeleteStateNoKey(t *testing.T) {
	ctx := goopstest.Context{
		Charm: DeleteState,
	}

	stateIn := goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if ctx.CharmErr != nil {
		t.Fatalf("expected no error, got: %v", ctx.CharmErr)
	}

	if _, exists := stateOut.StoredState["my-key"]; exists {
		t.Errorf("expected StoredState[my-key] to be deleted, but it still exists")
	}

	if len(stateOut.StoredState) != 0 {
		t.Errorf("expected StoredState to be empty, got %d items", len(stateOut.StoredState))
	}
}
