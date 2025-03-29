package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestActionGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"banana"`),
		Err:    nil,
	}

	result, err := commands.ActionGet(fakeRunner, "fruit")
	if err != nil {
		t.Fatalf("ActionGet returned an error: %v", err)
	}

	if result != "banana" {
		t.Fatalf("Expected %q, got %q", "banana", result)
	}

	if fakeRunner.Command != commands.ActionGetCommand {
		t.Errorf("Expected command %q, got %q", commands.ActionGetCommand, fakeRunner.Command)
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

func TestActionFail_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	err := commands.ActionFail(fakeRunner, "my failure message")
	if err != nil {
		t.Fatalf("ActionFail returned an error: %v", err)
	}
	if fakeRunner.Command != commands.ActionFailCommand {
		t.Errorf("Expected command %q, got %q", commands.ActionFailCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "my failure message" {
		t.Errorf("Expected argument %q, got %q", "my failure message", fakeRunner.Args[0])
	}
}

func TestActionSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}
	actionSetValues := map[string]string{
		"fruit": "banana",
		"color": "yellow",
	}
	err := commands.ActionSet(fakeRunner, actionSetValues)
	if err != nil {
		t.Fatalf("ActionSet returned an error: %v", err)
	}
	if fakeRunner.Command != commands.ActionSetCommand {
		t.Errorf("Expected command %q, got %q", commands.ActionSetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "fruit=banana" && fakeRunner.Args[1] != "fruit=banana" {
		t.Errorf("Expected argument %q, got %q", "fruit=banana", fakeRunner.Args[0])
	}
	if fakeRunner.Args[0] != "color=yellow" && fakeRunner.Args[1] != "color=yellow" {
		t.Errorf("Expected argument %q, got %q", "color=yellow", fakeRunner.Args[1])
	}

}
