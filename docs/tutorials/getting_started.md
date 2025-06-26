---
description: A hands-on introduction to goops for new users.
---

# Getting Started

In this tutorial, we will write, build, and deploy a Go charm using `goops`. You can expect to spend about 10 minutes completing this tutorial.

## Pre-requisites

To complete this tutorial, you will need a Ubuntu 24.04 machine with the following specifications:

- **Memory**: 8GB
- **CPU**: 4 cores
- **Disk**: 30GB

You will also need the following software installed:

- [Go](https://snapcraft.io/go) (version 1.24 or later)
- [Charmcraft](http://snapcraft.io/charmcraft) (version 3.4 or later)
- [Juju](http://snapcraft.io/juju) (version 3.6 or later)

## 1. Write the Go charm using `goops`

Create a new directory for your charm project and initialize a Go module:

```bash
mkdir my-charm
cd my-charm
go mod init my-charm
```

Create a `cmd/my-charm/main.go` file with the following content:

```go
package main

import (
	"fmt"
	"os"

	"github.com/gruyaume/goops"
)

func main() {
	env := goops.ReadEnv()

	goops.LogInfof("Hook name: %s", env.HookName)

	err := Configure()
	if err != nil {
		goops.LogErrorf("Error handling hook: %s", err.Error())
		os.Exit(1)
	}
}

type Config struct {
	Username string `json:"username"`
}

func Configure() error {
	isLeader, err := goops.IsLeader()
	if err != nil {
		return fmt.Errorf("could not check if unit is leader: %w", err)
	}

	if !isLeader {
		err := goops.SetUnitStatus(goops.StatusBlocked, "Unit is not leader")
		if err != nil {
			return fmt.Errorf("could not set unit status: %w", err)
		}
		return nil
	}

	goops.LogInfof("Unit is leader")

	var config Config
	err = goops.GetConfig(&config)
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}

	if config.Username == "" {
		err := goops.SetUnitStatus(goops.StatusBlocked, "Username is not set in config")
		if err != nil {
			return fmt.Errorf("could not set unit status: %w", err)
		}
		return nil
	}

	err = goops.SetUnitStatus(goops.StatusActive, fmt.Sprintf("Username is set to '%s'", config.Username))
	if err != nil {
		return fmt.Errorf("could not set unit status: %w", err)
	}

	return nil
}
```

Install go dependencies:

```bash
go mod tidy
```

## 2. Add the charm definition

Create a `charmcraft.yaml` file in the root of your project with the following content:

```yaml
name: my-charm
summary: Example Juju charm that uses `goops`
description: |
  Example Juju charm that uses `goops`

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
      bin/my-charm: dispatch

config:
  options:
    username:
      type: string
      default: gruyaume
      description: >
        Example configuration option for the charm.
```

## 3. Build the charm

Run the following command to build your charm:

```bash
charmcraft pack --verbose
```

This will create a `my-charm_amd64.charm` file in the current directory.

## 4. Deploy the charm

Create a new Juju model:

```bash
juju add-model demo
```

Deploy the charm to the model:

```bash
juju deploy ./my-charm_amd64.charm --config username=pizza
```

Check the status of the deployed charm:

```bash
juju status
```

You should see the unit status as `active` with the message "Username is set to 'pizza'".

```shell
guillaume@courge:~/my-charm$ juju status
Model  Controller  Cloud/Region  Version  SLA          Timestamp
demo   k8s-jun22   k8s-jun22     3.6.7    unsupported  07:56:18-04:00

App       Version  Status  Scale  Charm     Channel  Rev  Address        Exposed  Message
my-charm           active      1  my-charm             2  10.152.183.98  no       Username is set to 'pizza'

Unit         Workload  Agent  Address     Ports  Message
my-charm/0*  active    idle   10.1.0.112         Username is set to 'pizza'
```

!!! success
    Congratulations! You have successfully written, built, and deployed a charm using `goops`. You can now explore more features of `goops` and enhance your charm further.
