package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestActionFail_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.FailActionf("my failure message")
	if err != nil {
		t.Fatalf("ActionFail returned an error: %v", err)
	}

	if fakeRunner.Command != "action-fail" {
		t.Errorf("Expected command %q, got %q", "action-fail", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "my failure message" {
		t.Errorf("Expected argument %q, got %q", "my failure message", fakeRunner.Args[0])
	}
}

type MyActionParams struct {
	Fruit string `json:"fruit"`
	Color string `json:"color"`
}

func TestActionGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"fruit": "banana"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	actionParams := MyActionParams{}

	err := goops.GetActionParams(&actionParams)
	if err != nil {
		t.Fatalf("ActionGet returned an error: %v", err)
	}

	if actionParams.Fruit != "banana" {
		t.Fatalf("Expected %q, got %q", "banana", actionParams.Fruit)
	}

	if actionParams.Color != "" {
		t.Fatalf("Expected color to be empty, got %q", actionParams.Color)
	}

	if fakeRunner.Command != "action-get" {
		t.Errorf("Expected command %q, got %q", "action-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}

func TestActionLog_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.ActionLogf("my log message")
	if err != nil {
		t.Fatalf("ActionLog returned an error: %v", err)
	}

	if fakeRunner.Command != "action-log" {
		t.Errorf("Expected command %q, got %q", "action-log", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "my log message" {
		t.Errorf("Expected argument %q, got %q", "my log message", fakeRunner.Args[0])
	}
}

func TestActionSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(``),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	actionSetValues := map[string]string{
		"fruit": "banana",
		"color": "yellow",
	}

	err := goops.SetActionResults(actionSetValues)
	if err != nil {
		t.Fatalf("ActionSet returned an error: %v", err)
	}

	if fakeRunner.Command != "action-set" {
		t.Errorf("Expected command %q, got %q", "action-set", fakeRunner.Command)
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
