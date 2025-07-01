---
description: Add a configuration option to your Kubernetes charm using `goops`.
---

# 2. Add a configuration option

We will add a configuration option to our `myapp` charm that allows the user to set the port on which the application listens.

## 2.1 Update the Go charm

Open the `internal/charm/charm.go` file and update it to include a configuration option for the port:

```go

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
