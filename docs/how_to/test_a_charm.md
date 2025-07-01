---
description: Test a Juju charm using `goops` and `goopstest`.
---

# How-to test a Charm

## 1. Write the test cases

Create a `<charm name>_test.go` file in the same directory as your charm code. This file will contain the test cases for your charm using the `goopstest` package. Here we assume the charm name is `example`:

```go
package charm_test

import (
	"example/internal/charm"
	"testing"

	"github.com/gruyaume/goops/goopstest"
)

func TestConfigure(t *testing.T) {
	ctx := goopstest.NewContext(charm.Configure)

	stateIn := goopstest.State{
		Leader: true,
		Config: map[string]any{
			"username": "",
		},
	}

	stateOut := ctx.Run("install", stateIn)

	if ctx.CharmErr != nil {
		t.Fatalf("Configure failed: %v", ctx.CharmErr)
	}

	expectedStatus := goopstest.Status{
		Name:    goopstest.StatusBlocked,
		Message: "Username is not set in config",
	}
	if stateOut.UnitStatus != expectedStatus {
		t.Errorf("expected unit status %v, got %v", expectedStatus, stateOut.UnitStatus)
	}
}
```

!!! info
    Learn more about `goopstest`:

      - [Unit testing explanation](../explanation/unit_testing.md)
      - [goopstest API :octicons-link-external-24:](https://pkg.go.dev/github.com/gruyaume/goops/goopstest)

## 2. Run the tests

Run the tests using `go test ./... -v`: 

```shell
(venv) guillaume@courge:~/example$ go test ./... -v
?       example/cmd/example     [no test files]
=== RUN   TestConfigure
--- PASS: TestConfigure (0.00s)
PASS
ok      example/internal/charm  0.002s
```
