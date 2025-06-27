---
description: Charm Libraries reference.
---

# Charm Libraries

## Integrations

In Juju, [**integrations**](https://documentation.ubuntu.com/juju/3.6/reference/relation/index.html) are connections between charms, allowing them to share data through standardized interfaces. Common integrations exist for database access, TLS certificates, and much more. Integration specifications are defined centrally at [github.com/canonical/charm-relation-interfaces](https://github.com/canonical/charm-relation-interfaces).

## Charm Libraries

Charm libraries are reusable code packages that provide access to these integrations, allowing charms to integrate with other charms without needing to implement the integration logic themselves. Charm libraries for charm written with `goops` are maintained centrally at [github.com/gruyaume/charm-libraries](https://github.com/gruyaume/charm-libraries).

We maintain charm libraries for the following integrations:

- [`tls_certificates`](https://github.com/gruyaume/charm-libraries/tree/main/certificates): Securely request and manage TLS certificates.
- [`loki_push_api`](https://github.com/gruyaume/charm-libraries/tree/main/logging): Push logs to a Loki instance.
- [`prometheus_scrape`](https://github.com/gruyaume/charm-libraries/tree/main/prometheus): Send metrics related information to a Prometheus instance to allow scraping.
- [`tracing`](https://github.com/gruyaume/charm-libraries/tree/main/tracing): Receive tracing URLs from a tracing server.
