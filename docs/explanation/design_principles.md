---
description: Design principles for `goops`.
---

# Design principles

Goops provides access to the following Juju concepts:

- **Hook Commands**: `goops` exposes every Juju [hook commands](https://documentation.ubuntu.com/juju/3.6/reference/hook-command/list-of-hook-commands/), as a Go function.
- **Environment Variables**: `goops` provides access to every Juju-defined [environment variables](https://documentation.ubuntu.com/juju/3.6/reference/hook/#hook-execution).
- **Charm metadata**: `goops` provides access to the charm metadata as defined in `charmcraft.yaml`.
- **Pebble**: `goops` provides access to the Pebble API, allowing you to manage services and containers for Kubernetes charms.

## 1. Reliability

Charms are meant to be deployed at scale, in production environments, and in mission-critical applications. `goops` is designed to be reliable and predictable, ensuring that your charms behave consistently across different environments.

## 2. Simplicity

`goops` serves as a minimal mapping between Juju and Go constructs. Goops is not a framework; it does not impose charm design patterns. You call `goops`, it does not call you.
