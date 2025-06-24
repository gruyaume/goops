package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestSetUnitStatus_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.SetUnitStatus(goops.StatusActive)
	if err != nil {
		t.Fatalf("SetUnitStatus returned an error: %v", err)
	}

	if fakeRunner.Command != "status-set" {
		t.Errorf("Expected command %q, got %q", "status-set", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != string(goops.StatusActive) {
		t.Errorf("Expected argument %q, got %q", string(goops.StatusActive), fakeRunner.Args[0])
	}

	if fakeRunner.Output != nil {
		t.Errorf("Expected no output, got %q", string(fakeRunner.Output))
	}
}

func TestSetAppStatus_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.SetAppStatus(goops.StatusActive)
	if err != nil {
		t.Fatalf("SetAppStatus returned an error: %v", err)
	}

	if fakeRunner.Command != "status-set" {
		t.Errorf("Expected command %q, got %q", "status-set", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--application" {
		t.Errorf("Expected argument %q, got %q", "--application", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != string(goops.StatusActive) {
		t.Errorf("Expected argument %q, got %q", string(goops.StatusActive), fakeRunner.Args[1])
	}

	if fakeRunner.Output != nil {
		t.Errorf("Expected no output, got %q", string(fakeRunner.Output))
	}
}

func TestGetUnitStatus_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"status": "active", "message": "Unit is active"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	status, err := goops.GetUnitStatus()
	if err != nil {
		t.Fatalf("GetUnitStatus returned an error: %v", err)
	}

	if status.Name != goops.StatusActive {
		t.Errorf("Expected status %q, got %q", goops.StatusActive, status.Name)
	}

	if status.Message != "Unit is active" {
		t.Errorf("Expected message %q, got %q", "Unit is active", status.Message)
	}

	if fakeRunner.Command != "status-get" {
		t.Errorf("Expected command %q, got %q", "status-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--include-data" {
		t.Errorf("Expected argument %q, got %q", "--include-data", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}

func TestGetAppStatus_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"application-status":{"message":"Application is active","status":"active","status-data":{},"units":{"example/0":{"message":"","status":"unknown","status-data":{}},"example/1":{"message":"Application is active","status":"active","status-data":{}}}}}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	status, err := goops.GetAppStatus()
	if err != nil {
		t.Fatalf("GetAppStatus returned an error: %v", err)
	}

	if status.Name != goops.StatusActive {
		t.Errorf("Expected status %q, got %q", goops.StatusActive, status.Name)
	}

	if status.Message != "Application is active" {
		t.Errorf("Expected message %q, got %q", "Application is active", status.Message)
	}

	if status.Units["example/0"].Name != goops.StatusUnknown {
		t.Errorf("Expected unit status %q, got %q", goops.StatusUnknown, status.Units["example/0"].Name)
	}

	if status.Units["example/1"].Name != goops.StatusActive {
		t.Errorf("Expected unit status %q, got %q", goops.StatusActive, status.Units["example/1"].Name)
	}

	if fakeRunner.Command != "status-get" {
		t.Errorf("Expected command %q, got %q", "status-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 3 {
		t.Fatalf("Expected 3 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--application" {
		t.Errorf("Expected argument %q, got %q", "--application", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--include-data" {
		t.Errorf("Expected argument %q, got %q", "--include-data", fakeRunner.Args[1])
	}

	if fakeRunner.Args[2] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[2])
	}
}
