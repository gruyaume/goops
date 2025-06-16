package goops_test

import (
	"testing"

	"github.com/gruyaume/goops"
)

func TestResourceGet_Success(t *testing.T) {
	fakeRunner := &FakeRunner{
		Output: []byte(`/var/lib/juju/agents/unit-resources-example-0/resources/software/software.zip`),
		Err:    nil,
	}

	goops.SetCommandRunner(fakeRunner)

	result, err := goops.GetResource("software")
	if err != nil {
		t.Fatalf("ResourceGet returned an error: %v", err)
	}

	if result != "/var/lib/juju/agents/unit-resources-example-0/resources/software/software.zip" {
		t.Fatalf("Expected %q, got %q", "/var/lib/juju/agents/unit-resources-example-0/resources/software/software.zip", result)
	}

	if fakeRunner.Command != "resource-get" {
		t.Errorf("Expected command %q, got %q", "resource-get", fakeRunner.Command)
	}

	if len(fakeRunner.Args) != 1 {
		t.Fatalf("Expected 2 arguments, got %d", len(fakeRunner.Args))
	}

	if fakeRunner.Args[0] != "software" {
		t.Errorf("Expected argument %q, got %q", "software", fakeRunner.Args[0])
	}
}
