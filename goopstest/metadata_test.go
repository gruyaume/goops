package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func GetMetadata() error {
	meta, err := goops.ReadMetadata()
	if err != nil {
		return err
	}

	if meta.Name != "example" {
		return fmt.Errorf("expected metadata name to be 'example', got '%s'", meta.Name)
	}

	if meta.Description != "An example charm" {
		return fmt.Errorf("expected metadata description to be 'An example charm', got '%s'", meta.Description)
	}

	if len(meta.Containers) != 1 {
		return fmt.Errorf("expected metadata to contain one container")
	}

	if len(meta.Provides) != 1 {
		return fmt.Errorf("expected metadata to provide one interface")
	}

	return nil
}

func TestGetMetadata(t *testing.T) {
	ctx := goopstest.Context{
		Metadata: goops.Metadata{
			Name:        "example",
			Description: "An example charm",
			Containers: map[string]goops.Container{
				"example-container": {
					Resource: "example-image",
					Mounts:   []goops.Mount{},
				},
			},
			Provides: map[string]goops.Integration{
				"example-interface": {
					Interface: "example-interface",
				},
			},
		},
		Charm: GetMetadata,
	}

	stateIn := &goopstest.State{}

	_, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}
