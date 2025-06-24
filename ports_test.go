package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestOpenPortTCP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.OpenPort(80, "tcp")
	if err != nil {
		t.Fatalf("OpenPort returned an error: %v", err)
	}

	if fakeRunner.Command != "open-port" {
		t.Errorf("Expected command %q, got %q", "open-port", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "80/tcp" {
		t.Errorf("Expected argument %q, got %q", "80/tcp", fakeRunner.Args[0])
	}
}

func TestOpenPortUDP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.OpenPort(80, "udp")
	if err != nil {
		t.Fatalf("OpenPort returned an error: %v", err)
	}

	if fakeRunner.Command != "open-port" {
		t.Errorf("Expected command %q, got %q", "open-port", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "80/udp" {
		t.Errorf("Expected argument %q, got %q", "80/udp", fakeRunner.Args[0])
	}
}

func TestOpenPortICMP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.OpenPort(0, "icmp")
	if err != nil {
		t.Fatalf("OpenPort returned an error: %v", err)
	}

	if fakeRunner.Command != "open-port" {
		t.Errorf("Expected command %q, got %q", "open-port", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "icmp" {
		t.Errorf("Expected argument %q, got %q", "icmp", fakeRunner.Args[0])
	}
}

func TestOpenPortInvalidPort_Failure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.OpenPort(-1, "tcp")
	if err == nil {
		t.Fatalf("OpenPort did not return an error for invalid port")
	}

	if fakeRunner.Command != "" {
		t.Errorf("Expected no command to be run, but got %q", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 0 {
		t.Fatalf("Expected no arguments, got %d", len(fakeRunner.Args))
	}
}

func TestOpenPortInvalidProtocol_Failure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.OpenPort(80, "invalid")
	if err == nil {
		t.Fatalf("OpenPort did not return an error for invalid protocol")
	}

	if fakeRunner.Command != "" {
		t.Errorf("Expected no command to be run, but got %q", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 0 {
		t.Fatalf("Expected no arguments, got %d", len(fakeRunner.Args))
	}
}

func TestClosePortTCP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.ClosePort(80, "tcp")
	if err != nil {
		t.Fatalf("ClosePort returned an error: %v", err)
	}

	if fakeRunner.Command != "close-port" {
		t.Errorf("Expected command %q, got %q", "close-port", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "80/tcp" {
		t.Errorf("Expected argument %q, got %q", "80/tcp", fakeRunner.Args[0])
	}
}

func TestClosePortUDP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.ClosePort(80, "udp")
	if err != nil {
		t.Fatalf("ClosePort returned an error: %v", err)
	}

	if fakeRunner.Command != "close-port" {
		t.Errorf("Expected command %q, got %q", "close-port", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "80/udp" {
		t.Errorf("Expected argument %q, got %q", "80/udp", fakeRunner.Args[0])
	}
}

func TestClosePortInvalidPort_Failure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.ClosePort(-1, "tcp")
	if err == nil {
		t.Fatalf("ClosePort did not return an error for invalid port")
	}

	if fakeRunner.Command != "" {
		t.Errorf("Expected no command to be run, but got %q", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 0 {
		t.Fatalf("Expected no arguments, got %d", len(fakeRunner.Args))
	}
}

func TestClosePortInvalidProtocol_Failure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.ClosePort(80, "invalid")
	if err == nil {
		t.Fatalf("ClosePort did not return an error for invalid protocol")
	}

	if fakeRunner.Command != "" {
		t.Errorf("Expected no command to be run, but got %q", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 0 {
		t.Fatalf("Expected no arguments, got %d", len(fakeRunner.Args))
	}
}
