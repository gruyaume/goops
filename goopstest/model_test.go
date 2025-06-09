package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetModelInfo() error {
	env := goops.ReadEnv()

	if env.ModelName != "test-model" {
		return fmt.Errorf("expected model name 'test-model', got '%s'", env.ModelName)
	}

	if env.ModelUUID != "12345678-1234-5678-1234-567812345678" {
		return fmt.Errorf("expected model UUID '12345678-1234-5678-1234-567812345678', got '%s'", env.ModelUUID)
	}

	return nil
}

func TestModelName(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetModelInfo,
	}

	model := &goopstest.Model{
		Name: "test-model",
		UUID: "12345678-1234-5678-1234-567812345678",
	}

	stateIn := &goopstest.State{
		Model: model,
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if stateOut.Model.Name != "test-model" {
		t.Errorf("got Model.Name=%q, want %q", stateOut.Model.Name, "test-model")
	}

	if stateOut.Model.UUID != "12345678-1234-5678-1234-567812345678" {
		t.Errorf("got Model.UUID=%q, want %q", stateOut.Model.UUID, "12345678-1234-5678-1234-567812345678")
	}
}
