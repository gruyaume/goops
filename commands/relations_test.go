package commands_test

import (
	"testing"

	"github.com/gruyaume/go-operator/commands"
)

func TestRelationIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}

	result, err := commands.RelationIDs(fakeRunner, "tls-certificates")
	if err != nil {
		t.Fatalf("RelationIDs returned an error: %v", err)
	}

	expectedOutput := []string{
		"123",
		"456",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d relation IDs, got %d", len(expectedOutput), len(result))
	}
	for i, id := range result {
		if id != expectedOutput[i] {
			t.Errorf("Expected relation ID %q, got %q", expectedOutput[i], id)
		}
	}
	if fakeRunner.Command != commands.RelationIDsCommand {
		t.Errorf("Expected command %q, got %q", commands.RelationIDsCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 argument, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "tls-certificates" {
		t.Errorf("Expected argument %q, got %q", "tls-certificates", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}

func TestRelationGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"username":"user1","password":"pass1"}`),
		Err:    nil,
	}

	result, err := commands.RelationGet(fakeRunner, "certificates:0", "tls-certificates-requirer/0", false)
	if err != nil {
		t.Fatalf("RelationGet returned an error: %v", err)
	}

	expectedOutput := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d relation content keys, got %d", len(expectedOutput), len(result))
	}
	for k, v := range result {
		if v != expectedOutput[k] {
			t.Errorf("Expected relation content %q: %q, got %q", k, expectedOutput[k], v)
		}
	}
	if fakeRunner.Command != commands.RelationGetCommand {
		t.Errorf("Expected command %q, got %q", commands.RelationGetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "-r=certificates:0" {
		t.Errorf("Expected argument %q, got %q", "-r=certificates:0", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "-" {
		t.Errorf("Expected argument %q, got %q", "-", fakeRunner.Args[1])
	}
	if fakeRunner.Args[2] != "tls-certificates-requirer/0" {
		t.Errorf("Expected argument %q, got %q", "tls-certificates-requirer/0", fakeRunner.Args[2])
	}
	if fakeRunner.Args[3] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[3])
	}
}

func TestRelationList_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["tls-certificates-requirer/0", "tls-certificates-requirer/1"]`),
		Err:    nil,
	}

	result, err := commands.RelationList(fakeRunner, "certificates:0")
	if err != nil {
		t.Fatalf("RelationList returned an error: %v", err)
	}
	expectedOutput := []string{
		"tls-certificates-requirer/0",
		"tls-certificates-requirer/1",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d relation list items, got %d", len(expectedOutput), len(result))
	}
	for i, item := range result {
		if item != expectedOutput[i] {
			t.Errorf("Expected relation list item %q, got %q", expectedOutput[i], item)
		}
	}

	if fakeRunner.Args[0] != "-r=certificates:0" {
		t.Errorf("Expected argument %q, got %q", "-r=certificates:0", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}

func TestRelationSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	err := commands.RelationSet(fakeRunner, "certificates:0", true, map[string]string{"username": "user1", "password": "pass1"})
	if err != nil {
		t.Fatalf("RelationSet returned an error: %v", err)
	}

	if fakeRunner.Command != commands.RelationSetCommand {
		t.Errorf("Expected command %q, got %q", commands.RelationSetCommand, fakeRunner.Command)
	}
	if len(fakeRunner.Args) != 4 {
		t.Fatalf("Expected 4 arguments, got %d", len(fakeRunner.Args))
	}
	if fakeRunner.Args[0] != "-r=certificates:0" {
		t.Errorf("Expected argument %q, got %q", "-r=certificates:0", fakeRunner.Args[0])
	}
	if fakeRunner.Args[1] != "--app" {
		t.Errorf("Expected argument %q, got %q", "--app", fakeRunner.Args[1])
	}
	if fakeRunner.Args[2] != "username=user1" && fakeRunner.Args[2] != "password=pass1" {
		t.Errorf("Expected argument %q or %q, got %q", "username=user1", "password=pass1", fakeRunner.Args[2])
	}
	if fakeRunner.Args[3] != "username=user1" && fakeRunner.Args[3] != "password=pass1" {
		t.Errorf("Expected argument %q or %q, got %q", "username=user1", "password=pass1", fakeRunner.Args[3])
	}
}
