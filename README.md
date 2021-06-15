# Telepresence Demo

This repository holds code for a demonstration of [Telepresence](https://www.telepresence.io/). Below describes the scripts and commands to setup the demo and run through it. 

The demonstration has three components, all running inside of a Kubernetes cluster, that are not accessable from outside of the cluster (which is the point of Telepresence):

* [MongoDB](https://docs.mongodb.com/)
* [GRPC Client](./cmd/client/main.go)
* [GRPC Server](./cmd/server/main.go)


## The Applications

The demo applications are simple. They are [GRPC](https://grpc.io/docs/) apps are written in [Go](https://golang.org).

### The Server 
The server has the following workflows:

* connects to MongoDB
* seeds data on start up
* listens for connections from the client
* reads data from MongoDB and streams back to the client upon request
* logs the streams of data

### The Client
The client has the following workflow

* every second it requests data from the Server and logs the streamed response



*Note: The apps are meant to be quickly written simple single file apps. There is obviously a lot improvements to be made and production best practices implimented. They are just single files with a main package and they import protobuf. 

## Protobuff

Build the protobuff definition

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto
```

### Deploy MongoDB

Apply the yamls

```sh
kubectl apply -f infra/
```
