# Contributing

`goops` is an open-source project and we welcome contributions from the community. This document provides guidelines for contributing to the project. Contributions to `goops` can be made in the form of code, documentation, bug reports, feature requests, and feedback. We will judge contributions based on their quality, relevance, and alignment with the project's tenets.

## How-to

### Run tests

```shell
go test -short ./...     # unit tests
go vet ./...             # static analysis
golangci-lint run ./...  # linting
```

### Run integration tests

Pre-requisites:
- A Juju controller

```shell
go test ./...  # integration tests
```
