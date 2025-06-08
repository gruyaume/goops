package goopstest_test

import (
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func SetPorts() error {
	err := goops.SetPorts([]*goops.Port{
		{
			Port:     80,
			Protocol: "tcp",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func TestSetPorts(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetPorts,
	}

	stateIn := &goopstest.State{}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(stateOut.Ports))
	}

	if stateOut.Ports[0].Port != 80 {
		t.Errorf("Expected port 80, got %d", stateOut.Ports[0].Port)
	}

	if stateOut.Ports[0].Protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got '%s'", stateOut.Ports[0].Protocol)
	}
}

func TestSetPortsAlreadySet(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetPorts,
	}

	stateIn := &goopstest.State{
		Ports: []*goopstest.Port{
			{
				Port:     80,
				Protocol: "tcp",
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(stateOut.Ports))
	}

	if stateOut.Ports[0].Port != 80 {
		t.Errorf("Expected port 80, got %d", stateOut.Ports[0].Port)
	}

	if stateOut.Ports[0].Protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got '%s'", stateOut.Ports[0].Protocol)
	}
}

func TestSetPortsDifferentSet(t *testing.T) {
	ctx := goopstest.Context{
		Charm: SetPorts,
	}

	stateIn := &goopstest.State{
		Ports: []*goopstest.Port{
			{
				Port:     81,
				Protocol: "udp",
			},
		},
	}

	stateOut, err := ctx.Run("start", stateIn)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if len(stateOut.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(stateOut.Ports))
	}

	if stateOut.Ports[0].Port != 80 {
		t.Errorf("Expected port 80, got %d", stateOut.Ports[0].Port)
	}

	if stateOut.Ports[0].Protocol != "tcp" {
		t.Errorf("Expected protocol 'tcp', got '%s'", stateOut.Ports[0].Protocol)
	}
}
