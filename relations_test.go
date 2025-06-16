package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestRelationIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.GetRelationIDs("tls-certificates")
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

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.GetUnitRelationData("certificates:0", "tls-certificates-requirer/0")
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

func TestListRelationUnits_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["tls-certificates-requirer/0", "tls-certificates-requirer/1"]`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.ListRelationUnits("certificates:0")
	if err != nil {
		t.Fatalf("ListRelationUnits returned an error: %v", err)
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

func TestGetRelationApp_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`"provider"`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	app, err := goops.GetRelationApp("certificates:0")
	if err != nil {
		t.Fatalf("GetRelationApp returned an error: %v", err)
	}

	expectedApp := "provider"

	if app != expectedApp {
		t.Fatalf("Expected app %q, got %q", expectedApp, app)
	}

	if fakeRunner.Command != "relation-list" {
		t.Errorf("Expected command %q, got %q", "relation-list", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(fakeRunner.Args))
	}
}

func TestRelationSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: nil,
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.SetAppRelationData("certificates:0", map[string]string{
		"username": "user1",
		"password": "pass1",
	})
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

	goops.SetCommandRunner(fakeRunner)

	uuid, err := goops.GetRelationModel("certificates:0")
	if err != nil {
		t.Fatalf("RelationModelGet returned an error: %v", err)
	}

	if uuid != "e7ba04d1-b5f2-4769-8ae2-22e9119bca60" {
		t.Fatalf("Expected UUID %q, got %q", "e7ba04d1-b5f2-4769-8ae2-22e9119bca60", uuid)
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
