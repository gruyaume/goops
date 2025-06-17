package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ContainerCanConnect() error {
	pebble := goops.Pebble("example")

	_, err := pebble.SysInfo()
	if err != nil {
		return err
	}

	return nil
}

func TestContainerCantConnect(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerCanConnect,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: false,
			},
		},
	}

	_, err := ctx.Run("install", stateIn)
	if err == nil {
		t.Fatalf("Run should have returned an error, but got nil")
	}
}

func TestContainerCanConnect(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ContainerCanConnect,
	}

	stateIn := &goopstest.State{
		Containers: []*goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
			},
		},
	}

	_, err := ctx.Run("install", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}
