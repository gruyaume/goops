---
description: Add unit tests to your `goops` charm using `goopstest`.
---

# 2. Add unit tests using `goopstest`

In this section, we will write unit tests for the charm using `goopstest`. This step of the tutorial assumes you have completed the previous step.

## 2.1 Write the unit tests

### 2.1.1 Status is blocked given invalid port configuration

Here we will write a unit test that validates that the unit stauts is blocked given an invalid configuration. Create a new file `internal/charm/charm_test.go` with the following content:

```go
package charm_test

import (
	"myapp-k8s-operator/internal/charm"
	"testing"

	"github.com/gruyaume/goops/goopstest"
)

func TestGivenBadPortConfigWhenAnyEventThenStatusBlocked(t *testing.T) {
	ctx := goopstest.NewContext(
		charm.Configure,
	)
	stateIn := goopstest.State{
		Config: map[string]any{
			"port": 0, // Invalid port
		},
		Containers: []goopstest.Container{
			{
				Name:       "myapp",
				CanConnect: true,
			},
		},
	}

	stateOut := ctx.Run("update-status", stateIn)

	expectedStatus := goopstest.Status{
		Name:    goopstest.StatusBlocked,
		Message: "invalid config: port must be between 1 and 65535",
	}
	if stateOut.UnitStatus != expectedStatus {
		t.Errorf("expected status %v, got %v", expectedStatus, stateOut.UnitStatus)
	}
}
```

### 2.1.2 Status is active given valid configuration

Now, we will write a unit test that validates that the unit status is active given a valid configuration. Add the following test to the same file:

```go

func TestGivenValidConfigWhenAnyEventThenStatusActive(t *testing.T) {
	ctx := goopstest.NewContext(
		charm.Configure,
	)
	stateIn := goopstest.State{
		Config: map[string]any{
			"port": 8080,
		},
		Containers: []goopstest.Container{
			{
				Name:       "myapp",
				CanConnect: true,
			},
		},
	}

	stateOut := ctx.Run("update-status", stateIn)

	expectedStatus := goopstest.Status{
		Name:    goopstest.StatusActive,
		Message: "service is running on port 8080",
	}
	if stateOut.UnitStatus != expectedStatus {
		t.Errorf("expected status %v, got %v", expectedStatus, stateOut.UnitStatus)
	}
}
```

### 2.1.3 Pebble layer is added given valid configuration

Finally, we will write a unit test that validates that the Pebble layer is added when the configuration is valid. Add the following test to the same file:

```go
func TestGivenValidConfigWhenAnyEventThenPebbleLayerIsAdded(t *testing.T) {
	ctx := goopstest.NewContext(
		charm.Configure,
	)
	stateIn := goopstest.State{
		Config: map[string]any{
			"port": 8080, // Valid port
		},
		Containers: []goopstest.Container{
			{
				Name:       "myapp",
				CanConnect: true,
			},
		},
	}

	stateOut := ctx.Run("update-status", stateIn)

	got := stateOut.Containers[0].Layers["myapp"]

	want := goopstest.Layer{
		Summary:     "MyApp layer",
		Description: "pebble config layer for MyApp",
		Services: map[string]goopstest.Service{
			"myapp": {
				Summary:  "My App Service",
				Command:  "myapp -config /etc/myapp/config.yaml",
				Startup:  "enabled",
				Override: "replace",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected pebble layer (-want +got):\n%s", diff)
	}
}
```

Update the import statement at the top of the file to include the `cmp` package:

```go
import (
    "github.com/google/go-cmp/cmp"
)
``` 

Update the packages:

```bash
go mod tidy
```

## 2.2 Run the tests

Run the tests using the `go test` command:

```bash
go test -cover ./...
```

You should see output indicating that the tests passed:

```bash
guillaume@courge:~/code/myapp-k8s-operator$ go test -cover ./...
        myapp-k8s-operator/cmd/myapp-k8s-operator               coverage: 0.0% of statements
ok      myapp-k8s-operator/internal/charm       (cached)        coverage: 75.0% of statements
```
