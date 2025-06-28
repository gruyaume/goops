---
description: Build a `goops` charm.
---

# How-to build a Charm

## 1. Create a `charmcraft.yaml` file with the Go plugin

To build a Go charm, you need to create a `charmcraft.yaml` file in the root of your charm project. Use the `go` plugin to build your charm in `charmcraft.yaml`:

```yaml
parts:
  charm:
    source: .
    plugin: go
    build-snaps:
      - go
    organize:
      bin/<your-charm-name>: dispatch
```

Here replace `<your-charm-name>` with the name of your charm. The `dispatch` file will be the entry point for your charm.

!!! info
    For more information on the charmcraft charm definition, read the [official charmcraft documentation](https://canonical-charmcraft.readthedocs-hosted.com/stable/reference/files/charmcraft-yaml-file/).

## 2. Build the charm

Build the charm:

```shell
charmcraft pack --verbose
```
