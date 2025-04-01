package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestOpenPortTCP_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "tcp",
	}

	err := command.OpenPort(port)
	if err != nil {
		t.Fatalf("OpenPort returned an error: %v", err)
	}

	if fakeRunner.Command != commands.OpenPortCommand {
		t.Errorf("Expected command %q, got %q", commands.OpenPortCommand, fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "udp",
	}

	err := command.OpenPort(port)
	if err != nil {
		t.Fatalf("OpenPort returned an error: %v", err)
	}

	if fakeRunner.Command != commands.OpenPortCommand {
		t.Errorf("Expected command %q, got %q", commands.OpenPortCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "80/udp" {
		t.Errorf("Expected argument %q, got %q", "80/udp", fakeRunner.Args[0])
	}
}

func TestOpenPortInvalidPort_Failure(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     -1,
		Protocol: "tcp",
	}

	err := command.OpenPort(port)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "invalid",
	}

	err := command.OpenPort(port)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "tcp",
	}

	err := command.ClosePort(port)
	if err != nil {
		t.Fatalf("ClosePort returned an error: %v", err)
	}

	if fakeRunner.Command != commands.ClosePortCommand {
		t.Errorf("Expected command %q, got %q", commands.ClosePortCommand, fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "udp",
	}

	err := command.ClosePort(port)
	if err != nil {
		t.Fatalf("ClosePort returned an error: %v", err)
	}

	if fakeRunner.Command != commands.ClosePortCommand {
		t.Errorf("Expected command %q, got %q", commands.ClosePortCommand, fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     -1,
		Protocol: "tcp",
	}

	err := command.ClosePort(port)
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
	command := commands.Command{
		Runner: fakeRunner,
	}
	port := commands.Port{
		Port:     80,
		Protocol: "invalid",
	}

	err := command.ClosePort(port)
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
