---
description: Add a configuration option to your Kubernetes charm using `goops`.
---

# 2. Make the port configurable

We will add a configuration option to our `myapp` charm that allows the user to set the port on which the application listens. This step of the tutorial assumes you have completed the previous step.

## 2.1 Update the Go charm

Open the `internal/charm/charm.go` file, remove the hardcoded port, read the port from the Juju configuration options, and use it to write the configuration file. Here's the updated code:

```go
package charm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
	"gopkg.in/yaml.v3"
)

const (
	ConfigPath = "/etc/myapp/config.yaml"
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

type JujuConfig struct {
	Port int `json:"port"`
}

func Configure() error {
	c := JujuConfig{}

	err := goops.GetConfig(&c)
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	if c.Port < 1 || c.Port > 65535 {
		_ = goops.SetUnitStatus(goops.StatusBlocked, "invalid config: port must be between 1 and 65535")
		return nil
	}

	err = goops.SetPorts([]*goops.Port{
		{Port: c.Port, Protocol: goops.ProtocolTCP},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	pebble := goops.Pebble("myapp")

	_, err = pebble.SysInfo()
	if err != nil {
		_ = goops.SetUnitStatus(goops.StatusWaiting, "waiting for pebble to be ready")
		return nil
	}

	configFileChanged, err := syncConfig(pebble, c.Port)
	if err != nil {
		return fmt.Errorf("could not sync config: %w", err)
	}

	err = syncPebbleService(pebble)
	if err != nil {
		return fmt.Errorf("could not sync pebble service: %w", err)
	}

	if configFileChanged {
		_, err = pebble.Restart(
			&client.ServiceOptions{
				Names: []string{"myapp"},
			},
		)
		if err != nil {
			return fmt.Errorf("could not start pebble service: %w", err)
		}

		goops.LogInfof("Pebble service restarted")
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "service is running")

	return nil
}

type MyAppConfig struct {
	Port int `yaml:"port"`
}

func getExpectedConfig(port int) ([]byte, error) {
	myappConfig := MyAppConfig{
		Port: port,
	}

	b, err := yaml.Marshal(myappConfig)
	if err != nil {
		return nil, fmt.Errorf("could not marshal config to YAML: %w", err)
	}

	return b, nil
}

func syncConfig(pebble goops.PebbleClient, port int) (bool, error) {
	expectedContent, err := getExpectedConfig(port)
	if err != nil {
		return false, fmt.Errorf("could not get expected config: %w", err)
	}

	target := &bytes.Buffer{}

	err = pebble.Pull(&client.PullOptions{
		Path:   ConfigPath,
		Target: target,
	})
	if err != nil {
		goops.LogInfof("could not pull existing config from pebble: %v", err)
	}

	if target.String() == string(expectedContent) {
		goops.LogInfof("Config file is already up to date at %s", ConfigPath)
		return false, nil
	}

	err = pebble.Push(&client.PushOptions{
		Source: strings.NewReader(string(expectedContent)),
		Path:   ConfigPath,
	})
	if err != nil {
		return false, fmt.Errorf("could not push config to pebble: %w", err)
	}

	goops.LogInfof("Config file pushed to %s", ConfigPath)

	return true, nil
}

func syncPebbleService(pebble goops.PebbleClient) error {
	err := addPebbleLayer(pebble)
	if err != nil {
		return fmt.Errorf("could not add pebble layer: %w", err)
	}

	goops.LogInfof("Pebble layer created")

	_, err = pebble.Start(&client.ServiceOptions{
		Names: []string{"myapp"},
	})
	if err != nil {
		return fmt.Errorf("could not start pebble service: %w", err)
	}

	goops.LogInfof("Pebble service started")

	return nil
}

func addPebbleLayer(pebble goops.PebbleClient) error {
	layerData, err := yaml.Marshal(PebbleLayer{
		Summary:     "MyApp layer",
		Description: "pebble config layer for MyApp",
		Services: map[string]ServiceConfig{
			"myapp": {
				Override: "replace",
				Summary:  "My App Service",
				Command:  "myapp -config /etc/myapp/config.yaml",
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

## 2.2 Update the charm definition

Add a configuration option to the `charmcraft.yaml` file:

```yaml
...
config:
  options:
    port:
      type: int
      default: 8080
      description: >
        The port on which the application will listen.
```

## 2.3 Build the charm

Build the charm:

```bash
charmcraft pack --verbose
```

## 2.4 Deploy the charm

Deploy the charm to your Kubernetes cluster:

```bash
juju deploy ./myapp-k8s_amd64.charm --resource myapp-image=ghcr.io/gruyaume/myapp:latest
```

Open a web browser and navigate to the address of the `myapp-k8s` application on port 8080.

## 2.5 Change the configuration option

Change the configuration option to a different port:

```bash
juju config myapp-k8s port=8081
```

Open a web browser and navigate to the address of the `myapp-k8s` application on port 8081. You should see the same application running, but now on port 8081.
