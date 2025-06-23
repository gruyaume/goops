package goops_test

import (
	"testing"
	"time"

	"github.com/gruyaume/goops"
)

func TestSecretIDs_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`["123", "456"]`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.GetSecretIDs()
	if err != nil {
		t.Fatalf("SecretIDs returned an error: %v", err)
	}

	expectedOutput := []string{
		"123",
		"456",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d secret IDs, got %d", len(expectedOutput), len(result))
	}

	for i, id := range result {
		if id != expectedOutput[i] {
			t.Errorf("Expected secret ID %q, got %q", expectedOutput[i], id)
		}
	}

	if fakeRunner.Command != "secret-ids" {
		t.Errorf("Expected command %q, got %q", "secret-ids", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--format=json" {
		t.Errorf("Expected argument %q, got %q", "--format=json", fakeRunner.Args[0])
	}
}

func TestSecretGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"username":"user1","password":"pass1"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.GetSecretByLabel("my-label", false, true)
	if err != nil {
		t.Fatalf("SecretGet returned an error: %v", err)
	}

	expectedOutput := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	if len(result) != len(expectedOutput) {
		t.Fatalf("Expected %d secret content keys, got %d", len(expectedOutput), len(result))
	}

	for key, value := range result {
		if value != expectedOutput[key] {
			t.Errorf("Expected secret content %q, got %q", expectedOutput[key], value)
		}
	}

	if fakeRunner.Command != "secret-get" {
		t.Errorf("Expected command %q, got %q", "secret-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--label=my-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-label", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--refresh" {
		t.Errorf("Expected refresh arg %q, got %q", "--refresh", fakeRunner.Args[1])
	}

	if fakeRunner.Args[2] != "--format=json" {
		t.Errorf("Expected format arg %q, got %q", "--format=json", fakeRunner.Args[2])
	}
}

func TestSecretAdd_Success(t *testing.T) {
	content := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	description := "my secret"
	label := "my-label"

	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	expiry := time.Now().Add(24 * time.Hour)

	secretAddOptions := &goops.AddSecretOptions{
		Content:     content,
		Description: description,
		Expire:      time.Now().Add(24 * time.Hour),
		Label:       label,
		Rotate:      goops.RotateNever,
	}

	result, err := goops.AddSecret(secretAddOptions)
	if err != nil {
		t.Fatalf("SecretAdd returned an error: %v", err)
	}

	expectedOutput := `{"result":"success"}`
	if result != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, result)
	}

	if fakeRunner.Command != "secret-add" {
		t.Errorf("Expected command %q, got %q", "secret-add", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 6 {
		t.Fatalf("Expected 6 arguments, got %d", len(fakeRunner.Args))
	}

	contentArg := fakeRunner.Args[0]
	if fakeRunner.Args[0] != "username=user1" && fakeRunner.Args[1] != "username=user1" {
		t.Errorf("Expected content arg %q, got %q", "username=user1", contentArg)
	}

	if fakeRunner.Args[0] != "password=pass1" && fakeRunner.Args[1] != "password=pass1" {
		t.Errorf("Expected content arg %q, got %q", "password=pass1", contentArg)
	}

	if fakeRunner.Args[2] != "--description=my secret" {
		t.Errorf("Expected description arg %q, got %q", "--description=my secret", fakeRunner.Args[2])
	}

	if fakeRunner.Args[3] != "--label=my-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-label", fakeRunner.Args[3])
	}

	if fakeRunner.Args[4] != "--rotate=never" {
		t.Errorf("Expected rotate arg %q, got %q", "--rotate=never", fakeRunner.Args[4])
	}

	if fakeRunner.Args[5] != "--expire="+expiry.Format(time.RFC3339) {
		t.Errorf("Expected expire arg %q, got %q", "--expire="+expiry.Format(time.RFC3339), fakeRunner.Args[5])
	}
}

func TestSecretAdd_EmptyContent(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(""),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	secretAddOptions := &goops.AddSecretOptions{
		Description: "my secret",
		Expire:      time.Now().Add(24 * time.Hour),
		Label:       "my-label",
		Rotate:      goops.RotateNever,
	}

	_, err := goops.AddSecret(secretAddOptions)
	if err == nil {
		t.Error("Expected error when content is empty, but got nil")
	}
}

func TestSecretGrant_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.GrantSecretToRelation("123", "certificates:0")
	if err != nil {
		t.Fatalf("SecretGrant returned an error: %v", err)
	}

	if fakeRunner.Command != "secret-grant" {
		t.Fatalf("Expected command %q, got %q", "secret-grant", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "123" {
		t.Errorf("Expected ID arg %q, got %q", "123", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--relation=certificates:0" {
		t.Errorf("Expected secret ID arg %q, got %q", "--relation=certificates:0", fakeRunner.Args[1])
	}
}

func TestSecretInfoGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"cvng45vmp25c78cpk4u0":{"revision":1,"label":"active-ca-certificates","owner":"application","rotation":"never"}}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	secretInfo, err := goops.GetSecretInfoByLabel("my-secret-label")
	if err != nil {
		t.Fatalf("SecretInfoGet returned an error: %v", err)
	}

	expectedOutput := map[string]goops.SecretInfo{
		"cvng45vmp25c78cpk4u0": {
			Revision: 1,
			Label:    "active-ca-certificates",
			Owner:    "application",
			Rotation: "never",
		},
	}
	if len(secretInfo) != len(expectedOutput) {
		t.Fatalf("Expected %d secret info entries, got %d", len(expectedOutput), len(secretInfo))
	}

	if secretInfo["cvng45vmp25c78cpk4u0"].Revision != expectedOutput["cvng45vmp25c78cpk4u0"].Revision {
		t.Errorf("Expected revision %d, got %d", expectedOutput["cvng45vmp25c78cpk4u0"].Revision, secretInfo["cvng45vmp25c78cpk4u0"].Revision)
	}

	if secretInfo["cvng45vmp25c78cpk4u0"].Label != expectedOutput["cvng45vmp25c78cpk4u0"].Label {
		t.Errorf("Expected label %q, got %q", expectedOutput["cvng45vmp25c78cpk4u0"].Label, secretInfo["cvng45vmp25c78cpk4u0"].Label)
	}

	if secretInfo["cvng45vmp25c78cpk4u0"].Owner != expectedOutput["cvng45vmp25c78cpk4u0"].Owner {
		t.Errorf("Expected owner %q, got %q", expectedOutput["cvng45vmp25c78cpk4u0"].Owner, secretInfo["cvng45vmp25c78cpk4u0"].Owner)
	}

	if secretInfo["cvng45vmp25c78cpk4u0"].Rotation != expectedOutput["cvng45vmp25c78cpk4u0"].Rotation {
		t.Errorf("Expected rotation %q, got %q", expectedOutput["cvng45vmp25c78cpk4u0"].Rotation, secretInfo["cvng45vmp25c78cpk4u0"].Rotation)
	}

	if fakeRunner.Command != "secret-info-get" {
		t.Fatalf("Expected command %q, got %q", "secret-info-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "--label=my-secret-label" {
		t.Errorf("Expected label arg %q, got %q", "--label=my-secret-label", fakeRunner.Args[0])
	}

	if fakeRunner.Args[1] != "--format=json" {
		t.Errorf("Expected format arg %q, got %q", "--format=json", fakeRunner.Args[1])
	}
}

func TestSecretRemove_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.RemoveSecret("123")
	if err != nil {
		t.Fatalf("SecretRemove returned an error: %v", err)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "123" {
		t.Errorf("Expected ID arg %q, got %q", "123", fakeRunner.Args[0])
	}
}

func TestSecretRevoke_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	err := goops.RevokeSecret("123")
	if err != nil {
		t.Fatalf("SecretRevoke returned an error: %v", err)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 1 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "123" {
		t.Errorf("Expected ID arg %q, got %q", "123", fakeRunner.Args[0])
	}
}

func TestSecretSet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`{"result":"success"}`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	secretContent := map[string]string{
		"username": "user1",
		"password": "pass1",
	}
	expiry := time.Now().Add(24 * time.Hour)

	setSecretOpts := &goops.SetSecretOptions{
		ID:      "123",
		Content: secretContent,
		Expire:  expiry,
		Label:   "my-label",
		Rotate:  goops.RotateNever,
	}

	err := goops.SetSecret(setSecretOpts)
	if err != nil {
		t.Fatalf("couldn't set secret: %v", err)
	}

	if len(fakeRunner.Args) != 6 {
		t.Fatalf("Expected 6 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[1] != "username=user1" && fakeRunner.Args[2] != "username=user1" {
		t.Errorf("Expected ID arg %q", "username=user1")
	}

	if fakeRunner.Args[1] != "password=pass1" && fakeRunner.Args[2] != "password=pass1" {
		t.Errorf("Expected ID arg %q", "password=pass1")
	}

	if fakeRunner.Args[3] != "--label=my-label" {
		t.Errorf("Expected ID arg %q, got %q", "--label=my-label", fakeRunner.Args[3])
	}

	if fakeRunner.Args[4] != "--rotate=never" {
		t.Errorf("Expected ID arg %q, got %q", "--rotate=never", fakeRunner.Args[4])
	}

	if fakeRunner.Args[5] != "--expire="+expiry.Format(time.RFC3339) {
		t.Errorf("Expected ID arg %q, got %q", "--expire="+expiry.Format(time.RFC3339), fakeRunner.Args[5])
	}
}
