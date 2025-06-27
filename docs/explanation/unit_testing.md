---
description: Unit testing for `goops` charms.
---

# Unit Testing with `goopstest`

`goopstest` is a unit testing framework for `goops` charms. It allows you to simulate Juju environments and test your charm logic without needing a live Juju controller.

`goopstest` allows users to write unit tests in a "state-transition" style. Each test consists of:

- A Context and an initial state (Arrange)
- An event (Act)
- An output state (Assert)

## Examples

### A basic charm

Here's an example of a simple charm that uses `goops` to check if the unit is a leader and set its status accordingly:

```go
package charm

import (
	"github.com/gruyaume/goops"
)

func Configure() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return err
	}

	if !isLeader {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "Unit is not a leader")
		return nil
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "Charm is active")

	return nil
}
```

And here's the corresponding unit test written using `goopstest`:

```go
package charm_test

import (
	"testing"

	"github.com/gruyaume/goops/goopstest"
)

func TestCharm(t *testing.T) {
	// Arrange
	ctx := goopstest.NewContext(Configure)

	stateIn := goopstest.State{
		Leader: false,
	}

	// Act
	stateOut := ctx.Run("install", stateIn)

	// Assert
	expectedStatus := goopstest.Status{
		Name:    goopstest.StatusBlocked,
		Message: "Unit is not a leader",
	}
	if stateOut.UnitStatus != expectedStatus {
		t.Errorf("Expected unit status %v, got %v", expectedStatus, stateOut.UnitStatus)
	}
}
```

### A Kubernetes charm

Here's a Kubernetes charm example that uses `goops` to configure a Pebble service and start it:

```go
package charm

import (
	"fmt"
	"strings"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

func Configure() error {
	pebble := goops.Pebble("example")

	_, err := pebble.SysInfo()
	if err != nil {
		return fmt.Errorf("cannot connect to Pebble: %w", err)
	}

	err = pebble.Push(&client.PushOptions{
		Source: strings.NewReader(`# Example configuration file`),
		Path:   "/etc/config.yaml",
	})
	if err != nil {
		return fmt.Errorf("could not push file: %w", err)
	}

	layerData, err := yaml.Marshal(PebbleLayer{
		Summary:     "My service layer",
		Description: "This layer configures my service",
		Services: map[string]ServiceConfig{
			"my-service": {
				Startup:  "enabled",
				Override: "replace",
				Command:  "/bin/my-service --config /etc/my-service/config.yaml",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not marshal layer data to YAML: %w", err)
	}

	err = pebble.AddLayer(&client.AddLayerOptions{
		Combine:   true,
		Label:     "example-layer",
		LayerData: layerData,
	})
	if err != nil {
		return fmt.Errorf("could not add Pebble layer: %w", err)
	}

	_, err = pebble.Start(&client.ServiceOptions{
		Names: []string{"my-service"},
	})
	if err != nil {
		return fmt.Errorf("could not start Pebble service: %w", err)
	}

	return nil
}
```

And here's the corresponding unit test using `goopstest`:

```go
package charm_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops/goopstest"
)

func TestCharm(t *testing.T) {
	// Arrange
	ctx := goopstest.NewContext(Configure)

	dname, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(dname)

	stateIn := goopstest.State{
		Containers: []goopstest.Container{
			{
				Name:       "example",
				CanConnect: true,
				Mounts: map[string]goopstest.Mount{
					"config": {
						Location: "/etc/config.yaml",
						Source:   dname,
					},
				},
			},
		},
	}

	// Act
	stateOut := ctx.Run("install", stateIn)

	// Assert
	if len(stateOut.Containers) != 1 {
		t.Fatalf("Expected 1 container in stateOut, got %d", len(stateOut.Containers))
	}

	if len(stateOut.Containers[0].Layers) != 1 {
		t.Fatalf("Expected 1 Pebble layer in container, got %d", len(stateOut.Containers[0].Layers))
	}

	expectedLayer := goopstest.Layer{
		Summary:     "My service layer",
		Description: "This layer configures my service",
		Services: map[string]goopstest.Service{
			"my-service": {
				Startup:  "enabled",
				Override: "replace",
				Command:  "/bin/my-service --config /etc/my-service/config.yaml",
			},
		},
		LogTargets: map[string]*goopstest.LogTarget{},
	}

	actualLayer := stateOut.Containers[0].Layers["example-layer"]
	if !reflect.DeepEqual(actualLayer, expectedLayer) {
		t.Fatalf("Expected Pebble layer 'example-layer' to match expected configuration.\nExpected: %+v\nActual: %+v", expectedLayer, actualLayer)
	}

	if stateOut.Containers[0].ServiceStatuses["my-service"] != client.StatusActive {
		t.Errorf("Expected service 'my-service' to be active, got %s", stateOut.Containers[0].ServiceStatuses["my-service"])
	}

	content, err := os.ReadFile(dname + "/etc/config.yaml")
	if err != nil {
		t.Fatalf("Failed to read pushed file: %v", err)
	}

	if string(content) != "# Example configuration file" {
		t.Errorf("Expected file content '# Example configuration file', got '%s'", string(content))
	}
}
```

## Other resources

You can find more information about unit testing with `goopstest` in the following resources:

- [How-to: test a charm](../how_to/test_a_charm.md)
- [goopstest API Documentation :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops/goopstest)
