---
description: Design principles for `goops`.
---

# Design principles

`goops` is designed with a set of principles that guide its development. These principles are intended to ensure that `goops` remains focused on its core mission: developing reliable and simple charms in Go.

## 1. Reliability

Charms are meant to be deployed at scale, in production environments, and in mission-critical applications. `goops` is designed to be reliable and predictable, ensuring that your charms behave consistently across different environments.

## 2. Simplicity

`goops` serves as a minimal mapping between Juju and Go constructs. Goops is not a framework; it does not impose charm design patterns. You call `goops`, it does not call you.
