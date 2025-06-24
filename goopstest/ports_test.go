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

func ClosePort() error {
	err := goops.ClosePort(80, "udp")
	if err != nil {
		return err
	}

	return nil
}

func TestCloseUnOpenedPort(t *testing.T) {
	ctx := goopstest.Context{
		Charm: ClosePort,
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

	if ctx.CharmErr != nil {
		t.Errorf("Expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(stateOut.Ports))
	}

	if stateOut.Ports[0].Port != 81 {
		t.Errorf("Expected port 81, got %d", stateOut.Ports[0].Port)
	}

	if stateOut.Ports[0].Protocol != "udp" {
		t.Errorf("Expected protocol 'udp', got '%s'", stateOut.Ports[0].Protocol)
	}
}

func OpenPort() error {
	err := goops.OpenPort(81, "udp")
	if err != nil {
		return err
	}

	return nil
}

func TestOpenOpenedPort(t *testing.T) {
	ctx := goopstest.Context{
		Charm: OpenPort,
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

	if ctx.CharmErr != nil {
		t.Errorf("Expected no error, got %v", ctx.CharmErr)
	}

	if len(stateOut.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(stateOut.Ports))
	}

	if stateOut.Ports[0].Port != 81 {
		t.Errorf("Expected port 81, got %d", stateOut.Ports[0].Port)
	}

	if stateOut.Ports[0].Protocol != "udp" {
		t.Errorf("Expected protocol 'udp', got '%s'", stateOut.Ports[0].Protocol)
	}
}
