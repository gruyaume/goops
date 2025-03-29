# go-operator

** Write Juju charms in Go**

`go-operator` is a Go library for developing Juju charms.

> :construction: **Beta Notice**
> Go-operator is in beta. If you encounter any issues, please [report them here](https://github.com/gruyaume/go-operator/issues). 

## Try it now

**TO DO**

## Design principles

- **Reliability**: The most important attribute of a charm; therefore, it is the most important attribute of a charm framework.
- **Simplicity**: The library should be simple to use and understand via one-to-one mapping between Juju concepts and Go constructs.

## Why write a Charm in Go?

Reliability is the most important attribute of a charm. DevOps teams should focus their attention on operating their services, and the Charms they use should be reliable and predictable. The official Charm framework is written in Python, and while Python is great for getting started quickly, it is a poor long-term choice for production services. Contrary to popular belief, Python's biggest weakness is not its performance but its unreliability. Charms rely on the Python interpreter, which means that users may run charms in environments that significantly differ from the environment where developers built and tested them, leading to unexpected behavior. Python charms can't easily be ported from one system to another because of the differences in the Python interpreter and libraries. Go, on the other hand, is a compiled language that produces a single binary. Go charms are **reliable**, they **run on any system**, and they are **fast**.
