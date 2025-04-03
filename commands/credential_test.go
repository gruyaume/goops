package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestCredentialGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"key": "value"}`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	result, err := command.CredentialGet()
	if err != nil {
		t.Fatalf("CredentialGet returned an error: %v", err)
	}

	expected := map[string]string{"key": "value"}
	if result["key"] != expected["key"] {
		t.Fatalf("Expected %q, got %q", expected["key"], result["key"])
	}

	if fakeRunner.Command != commands.CredentialGetCommand {
		t.Errorf("Expected command %q, got %q", commands.CredentialGetCommand, fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}
