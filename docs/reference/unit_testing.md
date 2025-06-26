---
description: Unit testing for `goops` charms.
---

# Unit Testing

`goopstest` is a unit testing framework for `goops` charms. It allows you to simulate Juju environments and test your charm logic without needing a live Juju controller.

`goopstest` allows users to write unit tests in a "state-transition" style. Each test consists of:
- A Context and an initial state (Arrange)
- An event (Act)
- An output state (Assert)

You can refer to the [goopstest API documentation](https://pkg.go.dev/github.com/gruyaume/goops/goopstest) for more details.
