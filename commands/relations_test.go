package commands_test

import (
	"testing"

	"github.com/gruyaume/goops/commands"
)

func TestRelationIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	relationIDsOptions := &commands.RelationIDsOptions{
		Name: "tls-certificates",
	}

	result, err := command.RelationIDs(relationIDsOptions)
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

	if fakeRunner.Command != "relation-ids" {
		t.Errorf("Expected command %q, got %q", "relation-ids", fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}

	relationGetOptions := &commands.RelationGetOptions{
		ID:     "certificates:0",
		UnitID: "tls-certificates-requirer/0",
	}

	result, err := command.RelationGet(relationGetOptions)
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

	if fakeRunner.Command != "relation-get" {
		t.Errorf("Expected command %q, got %q", "relation-get", fakeRunner.Command)
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
	command := commands.Command{
		Runner: fakeRunner,
	}

	relationListOptions := &commands.RelationListOptions{
		ID: "certificates:0",
	}

	result, err := command.RelationList(relationListOptions)
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

	if fakeRunner.Command != "relation-list" {
		t.Errorf("Expected command %q, got %q", "relation-list", fakeRunner.Command)
	}
}

func TestRelationSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	relationSetOptions := &commands.RelationSetOptions{
		ID:  "certificates:0",
		App: true,
		Data: map[string]string{
			"username": "user1",
			"password": "pass1",
		},
	}

	err := command.RelationSet(relationSetOptions)
	if err != nil {
		t.Fatalf("RelationSet returned an error: %v", err)
	}

	if fakeRunner.Command != "relation-set" {
		t.Errorf("Expected command %q, got %q", "relation-set", fakeRunner.Command)
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

func TestRelationModelGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"uuid":"e7ba04d1-b5f2-4769-8ae2-22e9119bca60"}`),
		Err:    nil,
	}
	command := commands.Command{
		Runner: fakeRunner,
	}

	relationModelGetOptions := &commands.RelationModelGetOptions{
		ID: "certificates:0",
	}

	result, err := command.RelationModelGet(relationModelGetOptions)
	if err != nil {
		t.Fatalf("RelationModelGet returned an error: %v", err)
	}

	expectedOutput := commands.RelationModel{
		UUID: "e7ba04d1-b5f2-4769-8ae2-22e9119bca60",
	}

	if result.UUID != expectedOutput.UUID {
		t.Fatalf("Expected UUID %q, got %q", expectedOutput.UUID, result.UUID)
	}

	if fakeRunner.Command != "relation-model-get" {
		t.Errorf("Expected command %q, got %q", "relation-model-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "-r=certificates:0" {
		t.Errorf("Expected argument %q, got %q", "-r=certificates:0", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}
