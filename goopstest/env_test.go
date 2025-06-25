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

func GetUnitName() error {
	env := goops.ReadEnv()

	if env.UnitName != "blou/0" {
		return fmt.Errorf("expected unit name 'blou/0', got '%s'", env.UnitName)
	}

	return nil
}

func GetJujuVersion() error {
	env := goops.ReadEnv()

	if env.Version != "1.2.3" {
		return fmt.Errorf("expected Juju version '1.2.3', got '%s'", env.Version)
	}

	return nil
}

func TestGetModelInfo(t *testing.T) {
	ctx := goopstest.Context{
		Charm: GetModelInfo,
	}

	model := goopstest.Model{
		Name: "test-model",
		UUID: "12345678-1234-5678-1234-567812345678",
	}

	stateIn := goopstest.State{
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

func TestGetUnitName(t *testing.T) {
	ctx := goopstest.Context{
		Charm:   GetUnitName,
		AppName: "blou",
		UnitID:  "blou/0",
	}

	stateIn := goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}

func TestGetJujuVersion(t *testing.T) {
	ctx := goopstest.Context{
		Charm:       GetJujuVersion,
		JujuVersion: "1.2.3",
	}

	stateIn := goopstest.State{}

	_, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}
