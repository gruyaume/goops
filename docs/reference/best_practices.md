---
description: Best practices for writing charms with `goops`.
---

# Charm Development best Practices

This document outlines best practices for writing charms using the `goops` framework. Following these practices will help ensure that your charm is robust.

## Write idempotent charm code

Charm code should be thought of as a reconciliation loop that applies the necessary changes to ensure that the charm's state matches the desired state. Charm code should be idempotent, meaning that running the code multiple times should not change the state of the charm if it is already in the desired state.

## Use Charm Libraries for managing relation data

If you need to read or write to relation data, use the appropriate [Charm Library](../reference/charm_libraries.md). If the library does not exist, consider creating it.

## Use `goopstest` for unit testing

Use `goopstest` to write unit tests for your charms in a state-transition style.

## Be wary of state

Be wary of managing state in your charm code. Maintaining state becomes increasingly complex the longer the charm is deployed as users upgrade it and the charm code evolves.

State includes:
- Stored State
- Secrets
- Relation data
- Files on disk

## Write clear, idiomatic Go code

Write clear, idiomatic Go code that is easy to read and understand. Learn more about Go best practices in [the Effective Go guide](https://go.dev/doc/effective_go) and [Google's Go Style Guide](https://google.github.io/styleguide/go/guide).
