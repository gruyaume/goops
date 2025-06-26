---
description: Charm Libraries reference.
---

# Charm Libraries

In Juju, [**integrations**](https://documentation.ubuntu.com/juju/3.6/reference/relation/index.html) are connections between charms, allowing them to share data through standardized interfaces. Common integrations exist for database access, TLS certificates, and much more. Integration specifications are defined centrally at [github.com/canonical/charm-relation-interfaces](https://github.com/canonical/charm-relation-interfaces).

Charm libraries are reusable code packages that provide access to these integrations, allowing charms to integrate with other charms without needing to implement the integration logic themselves. Charm libraries for charm written with `goops` are maintained centrally at [github.com/gruyaume/charm-libraries](https://github.com/gruyaume/charm-libraries).
