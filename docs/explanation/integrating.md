---
description: Integrating `goops` charms with other charms.
---

# Integrating `goops` charms

## Integrations

In Juju, [**integrations**](https://documentation.ubuntu.com/juju/3.6/reference/relation/index.html) are connections between charms, allowing them to share data through standardized interfaces. Common integrations exist for database access, TLS certificates, and much more. Integration specifications are defined centrally at [github.com/canonical/charm-relation-interfaces](https://github.com/canonical/charm-relation-interfaces).

## goops charms integrate with other charms

`goops` charms can integrate with other charms using the same integration specifications. We maintain a set of [Charm Libraries](../reference/charm_libraries.md) for commonly used interfaces, allowing `goops` charms to integrate with many charms in the ecosystem, whether they are written in Python, Go, or any other language.
