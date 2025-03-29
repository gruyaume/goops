package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestConfigGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"banana"`),
		Err:    nil,
	}

	result, err := commands.ConfigGet(fakeRunner, "fruit")
	if err != nil {
		t.Fatalf("ConfigGet returned an error: %v", err)
	}
	if result != "banana" {
		t.Fatalf("Expected %q, got %q", "banana", result)
	}

	if fakeRunner.Command != commands.ConfigGetCommand {
		t.Errorf("Expected command %q, got %q", commands.ConfigGetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "fruit" {
		t.Errorf("Expected argument %q, got %q", "fruit", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}
