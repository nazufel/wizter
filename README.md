# Wizter

This repository holds code for various demonstrations. The list of demonstration steps and discussion can be found below in [Included Demos](#included-demos). These demonstrations use a real containerized applications with a micro services architecture and are deployed to Kubernetes. The applications components are:

* [gRPC Client](./cmd/client/main.go) - client side application for the demonstration written in Go
* [gRPC Server](./cmd/server/main.go) - service side application for the demonstration written in Go
* [MongoDB](https://docs.mongodb.com/) - database for storing state of the application

## The Applications

The demo applications are simple. They are [GRPC](https://grpc.io/docs/) apps are written in [Go](https://golang.org).

The Wizter application is a wizards listing service as defined in the [wizard.proto](./wizard/wizard.proto) file. The [client](./cmd/client/main.go) requests a list of wizards from the [server](./cmd/server/main.go). The server retrieves a list of wizards from the database and [streams](https://grpc.io/docs/what-is-grpc/core-concepts/#server-streaming-rpc) them back to the client.

*Note: The apps are meant to be a quickly written simple single file. There are obviously a lot improvements to be made and production best practices implemented. They are just single files with a main package and they import protobuf. The goal of this repo is not to write the prettiest or best tested Go programs, but to have a running application to use for the included demo steps.

## Included Demos

* [Bazel](./demos/07-29-2021-bazel.md)
* [Telepresene](./demos/06-11-2021-telepresence.md)


## Contributing

This repository is not accepting contributions. 