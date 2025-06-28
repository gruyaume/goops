---
description: Manage workloads on Kubernetes with `goops` charms.
---

# How-to manage workloads on Kubernetes

This guide covers how to manage workloads on Kubernetes using `goops`. Workloads are services that run in sidecar containers, next to the charm container. The charm uses Pebble to manage these workloads. [Pebble](https://github.com/canonical/pebble) is a lightweight Linux service manager that allows the charm to manage services, files, and health checks in the application's container.

## 1. Declare containers

Declare the containers in you charm's `charmcraft.yaml` file. For example:

```yaml
containers:
  myapp:
    resource: myapp-image

resources:
  myapp-image:
    type: oci-image
    description: OCI image for my application
```

!!! note
    For more information on the `charmcraft.yaml` charm definition, read the [official charmcraft documentation](https://canonical-charmcraft.readthedocs-hosted.com/stable/reference/files/charmcraft-yaml-file/).

## 2. Manage workloads using `goops.Pebble`

You can manage workloads using the `goops.Pebble` API. In the following example, we initialize a Pebble client for the `myapp` container, push a configuration file, create a Pebble layer, and start the service.

```go
package charm

import (
	"fmt"
	"strings"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

const (
	ConfigPath = "/etc/config.yaml"
)

type ServiceConfig struct {
	Override string `yaml:"override"`
	Summary  string `yaml:"summary"`
	Command  string `yaml:"command"`
	Startup  string `yaml:"startup"`
}

type PebbleLayer struct {
	Summary     string                   `yaml:"summary"`
	Description string                   `yaml:"description"`
	Services    map[string]ServiceConfig `yaml:"services"`
}

type PebblePlan struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

func Configure() error {
	pebble := goops.Pebble("myapp")

	_, err := pebble.SysInfo()
	if err != nil {
		return fmt.Errorf("could not connect to pebble: %w", err)
	}

	err = syncConfig(pebble)
	if err != nil {
		return fmt.Errorf("could not sync config: %w", err)
	}

	err = syncPebbleService(pebble)
	if err != nil {
		return fmt.Errorf("could not sync pebble service: %w", err)
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "service is running")

	return nil
}

func syncConfig(pebble goops.PebbleClient) error {
	expectedConfig := "Example configuration file for MyApp"

	err := pebble.Push(&client.PushOptions{
		Source: strings.NewReader(expectedConfig),
		Path:   ConfigPath,
	})
	if err != nil {
		return fmt.Errorf("could not push file: %w", err)
	}

	return nil
}

func syncPebbleService(pebble goops.PebbleClient) error {
	if !pebbleLayerCreated(pebble) {
		goops.LogInfof("Pebble layer not created")

		err := addPebbleLayer(pebble)
		if err != nil {
			return fmt.Errorf("could not add pebble layer: %w", err)
		}

		goops.LogInfof("Pebble layer created")
	}

	_, err := pebble.Start(&client.ServiceOptions{
		Names: []string{"notary"},
	})
	if err != nil {
		return fmt.Errorf("could not start pebble service: %w", err)
	}

	goops.LogInfof("Pebble service started")

	return nil
}

func pebbleLayerCreated(pebble goops.PebbleClient) bool {
	dataBytes, err := pebble.PlanBytes(nil)
	if err != nil {
		return false
	}

	var plan PebblePlan

	err = yaml.Unmarshal(dataBytes, &plan)
	if err != nil {
		return false
	}

	service, exists := plan.Services["myapp"]
	if !exists {
		return false
	}

	if service.Command != "myapp --config "+ConfigPath {
		return false
	}

	return true
}

func addPebbleLayer(pebble goops.PebbleClient) error {
	layerData, err := yaml.Marshal(PebbleLayer{
		Summary:     "MyApp layer",
		Description: "pebble config layer for MyApp",
		Services: map[string]ServiceConfig{
			"myapp": {
				Override: "replace",
				Summary:  "My App Service",
				Command:  "myapp --config " + ConfigPath,
				Startup:  "enabled",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not marshal layer data to YAML: %w", err)
	}

	err = pebble.AddLayer(&client.AddLayerOptions{
		Combine:   true,
		Label:     "myapp",
		LayerData: layerData,
	})
	if err != nil {
		return fmt.Errorf("could not add pebble layer: %w", err)
	}

	return nil
}
```

!!! info
    Learn more about workload management in Kubernetes charms:

    - [Pebble documentation :octicons-link-external-24:](https://documentation.ubuntu.com/pebble/)
    - [goops API reference :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops)
