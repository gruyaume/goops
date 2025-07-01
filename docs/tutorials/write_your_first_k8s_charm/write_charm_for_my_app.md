---
description: Tutorial to write a Kubernetes charm for `myapp` using `goops`.
---

# 1. Write a charm for `myapp`

We will write a Kubernetes charm for an application named [myapp](https://github.com/gruyaume/myapp). This simple web application requires a configuration file that contains the port on which it listens.

## 1.1. Write the Go charm using `goops`

Create a new directory for your charm project and initialize a Go module:

```bash
mkdir myapp-k8s-operator
cd myapp-k8s-operator
go mod init myapp-k8s-operator
```

Create a `cmd/myapp-k8s-operator/main.go` file with the following content:

```go
package main

import (
	"myapp-k8s-operator/internal/charm"
	"os"

	"github.com/gruyaume/goops"
)

func main() {
	err := charm.Configure()
	if err != nil {
		goops.LogErrorf("Failed to configure charm: %v", err)
		os.Exit(1)
	}
}
```

Create a `internal/charm/charm.go` file with the following content:

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
	Port       = 8080
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

func Configure() error {
	pebble := goops.Pebble("myapp")

	err := goops.SetPorts([]*goops.Port{
		{Port: Port, Protocol: goops.ProtocolTCP},
	})
	if err != nil {
		return fmt.Errorf("could not set ports: %w", err)
	}

	err = syncConfig(pebble)
	if err != nil {
		return fmt.Errorf("could not sync config: %w", err)
	}

	_, err = pebble.SysInfo()
	if err != nil {
		return fmt.Errorf("could not connect to pebble: %w", err)
	}

	err = syncPebbleService(pebble)
	if err != nil {
		return fmt.Errorf("could not sync pebble service: %w", err)
	}

	_ = goops.SetUnitStatus(goops.StatusActive, "service is running")

	return nil
}

type MyAppConfig struct {
	Port int `yaml:"port"`
}

func getExpectedConfig() ([]byte, error) {
	myappConfig := MyAppConfig{
		Port: Port,
	}

	b, err := yaml.Marshal(myappConfig)
	if err != nil {
		return nil, fmt.Errorf("could not marshal config to YAML: %w", err)
	}

	return b, nil
}

func syncConfig(pebble goops.PebbleClient) error {
	content, err := getExpectedConfig()
	if err != nil {
		return fmt.Errorf("could not get expected config: %w", err)
	}

	source := strings.NewReader(string(content))

	err = pebble.Push(&client.PushOptions{
		Source: source,
		Path:   ConfigPath,
	})
	if err != nil {
		return fmt.Errorf("could not push config to pebble: %w", err)
	}

	goops.LogInfof("Config file pushed to %s", ConfigPath)

	return nil
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

Install the go dependencies:

```bash
go mod tidy
```

## 1.2. Add the charm definition

Create a `charmcraft.yaml` file in the root of your project with the following content:

```yaml
name: myapp-k8s
summary: A Kubernetes charm for `myapp`
description: |
  A Kubernetes charm for `myapp`.

type: charm
base: ubuntu@24.04
build-base: ubuntu@24.04
platforms:
  amd64:

parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/myapp-k8s-operator: dispatch

containers:
  myapp:
    resource: myapp-image
    mounts:
    - storage: config
      location: /etc/myapp

storage:
  config:
    type: filesystem
    minimum-size: 5M

resources:
  myapp-image:
    type: oci-image
    description: OCI image for myapp
    upstream-source: ghcr.io/gruyaume/myapp:v0.0.1
```

## 1.3. Build the charm

Build the charm using `charmcraft`:

```bash
charmcraft pack --verbose
```

This will create a `myapp-k8s_amd64.charm` file in the current directory.

## 1.4. Deploy the charm

Create a new Juju model:

```bash
juju add-model demo
```

Deploy the charm to the model:

```bash
juju deploy ./myapp-k8s_amd64.charm --resource myapp-image=ghcr.io/gruyaume/myapp:latest
```

Verify that the charm is running:

```bash
juju status
```

You should see the `myapp-k8s` application in the status output, indicating that it is active and running.

```shell
Model  Controller  Cloud/Region  Version  SLA          Timestamp
demo   k8s-jul1    k8s-jul1      3.6.7    unsupported  09:47:22-04:00

App        Version  Status  Scale  Charm      Channel  Rev  Address         Exposed  Message
myapp-k8s           active      1  myapp-k8s             2  10.152.183.113  no       service is running

Unit          Workload  Agent  Address    Ports  Message
myapp-k8s/0*  active    idle   10.1.0.95         service is running
```

## 1.5. Access the application

Open a web browser and navigate to the address of the `myapp-k8s` application. Here this address is `http://10.1.0.95:8080`, replace the IP address with the one shown in the unit address in the `juju status` output. You should see a page displaying `MyApp, "/"`.
