package goopstest_test

import (
	"fmt"
	"testing"

	"github.com/gruyaume/goops"
	"github.com/gruyaume/goops/goopstest"
)

func ResourceGet() error {
	resource, err := goops.GetResource("my-resource")
	if err != nil {
		return err
	}

	if resource != "/var/lib/juju/agents/unit-example-0/resources/somefile.txt" {
		return fmt.Errorf("expected resource path to be '/var/lib/juju/agents/unit-example-0/resources/somefile.txt', got %s", resource)
	}

	return nil
}

func TestResourceGet(t *testing.T) {
	ctx := goopstest.NewContext(ResourceGet, goopstest.WithAppName("example"), goopstest.WithMetadata(
		goopstest.Metadata{
			Resources: map[string]goopstest.ResourceMeta{
				"my-resource": {
					Description: "A test resource",
					Type:        "file",
					Filename:    "somefile.txt",
				},
			},
		},
	))

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Charm returned an error: %v", ctx.CharmErr)
	}
}

func TestResourceGetDoesntExist(t *testing.T) {
	ctx := goopstest.NewContext(ResourceGet)

	stateIn := goopstest.State{}

	_ = ctx.Run("start", stateIn)

	if ctx.CharmErr == nil {
		t.Fatalf("Expected charm to return an error for non-existent resource, but got nil")
	}

	expectedErr := "command resource-get failed: ERROR could not download resource: HTTP request failed: Get https://1.2.3.4:17070/model/7bc47acd-4a48-4d11-8f52-3c44656bcb94/units/unit-example-0/resources/\"my-resource\": resource#example/\"my-resource\" not found"
	if ctx.CharmErr.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, ctx.CharmErr.Error())
	}
}
