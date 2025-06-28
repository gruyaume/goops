---
description: Example Juju charms that use `goops`.
---

# Example charms

The following charms use `goops` and can be used as reference implementations:

- [**Certificates**](https://github.com/gruyaume/certificates-operator): A charm for provisioning TLS certificates using the `tls-certificates` integration.
- [**Notary K8s**](https://github.com/gruyaume/notary-k8s-operator): A Kubernetes charm for [Notary](https://github.com/canonical/notary), a TLS certificate authority for enterprise applications. It works on both Kubernetes and machine models.
- [**LEGO**](https://github.com/gruyaume/lego-operator): A charm for managing Let's Encrypt certificates using the [LEGO](https://github.com/go-acme/lego) client. It works on both Kubernetes and machine models.
- [**Core K8s**](https://github.com/ellanetworks/core-k8s-operator): A Kubernetes charm for operating [Ella Core](https://docs.ellanetworks.com/), a 5G core network.

Feel free to explore these charms for practical examples of how to use `goops` in your own charms.
