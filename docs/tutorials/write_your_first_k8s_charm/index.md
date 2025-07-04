# Write your first K8s charm

In this tutorial, we will write a Kubernetes charm in Go for a web application using `goops`. The application we will charm is called [myapp](https://github.com/gruyaume/myapp), a simple web application that displays `"MyApp, '/'"`. This tutorial will take about 30 minutes to complete and you will how to use `goops` to manage Pebble services, relations, and configurations. You will also build your charm and deploy it to a Kubernetes cluster.

- [1. Write a charm for `myapp`](write_charm_for_my_app.md)
- [2. Make the port configurable](make_port_configurable.md)
- [2. Integrate with Loki](integrate_with_loki.md)

At any moment, you can refer to the [MyApp K8s Operator GitHub repository](https://github.com/gruyaume/myapp-k8s-operator) for the complete code of the charm we will write in this tutorial.
