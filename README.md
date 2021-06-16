# Telepresence Demo

This repository holds code for a demonstration of [Telepresence](https://www.telepresence.io/). Below describes the scripts and commands to setup the demo and run through it. 

The demonstration has three components, all running inside of a Kubernetes cluster, that are not accessible from outside of the cluster (which is the point of Telepresence):

* [MongoDB](https://docs.mongodb.com/)
* [GRPC Client](./cmd/client/main.go)
* [GRPC Server](./cmd/server/main.go)


## The Applications

The demo applications are simple. They are [GRPC](https://grpc.io/docs/) apps are written in [Go](https://golang.org).


The application is a wizards listing service as defined in the [wizards.proto](./wizards/wizards.proto) file. The client requests a list of wizards from the server. The server retrieves a list of wizards from the database and [streams](https://grpc.io/docs/what-is-grpc/core-concepts/#server-streaming-rpc) them back to the client.

*Note: The apps are meant to be a quickly written simple single file. There are obviously a lot improvements to be made and production best practices implemented. They are just single files with a main package and they import protobuf. The goal of this repo is not to write the prettiest or best tested Go programs, but to demo Telepresence and how it can improve development workflows.

## Non-Application Dependencies

* [Telepresence CLI](https://www.telepresence.io/docs/latest/install/)
* [Kubernetes](https://kubernetes.io/) (using [Kind](https://kind.sigs.k8s.io/))
* [GNU Make](https://www.gnu.org/software/make/)

## Begin the Demo
Here are the steps for running the demo.
### Deploy the Application

The demo begins with deploying the first version of the wizards listing application to a Kubernetes cluster. This demo will use [Kind](https://kind.sigs.k8s.io/), but any Kubernetes cluster will do where the user has administrative access to deploy services to.

#### Create a Kind Cluster

```sh
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind
```

Make sure the `kind-kind` cluster is set to the current context in [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).

```sh
kubectl config current-context

kind-kind
```

Check cluster access

```sh
kubectl get all --all-namespaces

NAMESPACE            NAME                                             READY   STATUS    RESTARTS   AGE
kube-system          pod/coredns-558bd4d5db-4kpp2                     1/1     Running   0          3m23s
kube-system          pod/coredns-558bd4d5db-gprnq                     1/1     Running   0          3m23s
kube-system          pod/etcd-kind-control-plane                      1/1     Running   0          3m25s
kube-system          pod/kindnet-bntcb                                1/1     Running   0          3m23s
kube-system          pod/kube-apiserver-kind-control-plane            1/1     Running   0          3m25s
kube-system          pod/kube-controller-manager-kind-control-plane   1/1     Running   0          3m25s
kube-system          pod/kube-proxy-4bm4n                             1/1     Running   0          3m23s
kube-system          pod/kube-scheduler-kind-control-plane            1/1     Running   0          3m25s
local-path-storage   pod/local-path-provisioner-547f784dff-5j89p      1/1     Running   0          3m23s

NAMESPACE     NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
default       service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP                  3m38s
kube-system   service/kube-dns     ClusterIP   10.96.0.10   <none>        53/UDP,53/TCP,9153/TCP   3m36s

NAMESPACE     NAME                        DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR            AGE
kube-system   daemonset.apps/kindnet      1         1         1       1            1           <none>                   3m35s
kube-system   daemonset.apps/kube-proxy   1         1         1       1            1           kubernetes.io/os=linux   3m36s

NAMESPACE            NAME                                     READY   UP-TO-DATE   AVAILABLE   AGE
kube-system          deployment.apps/coredns                  2/2     2            2           3m36s
local-path-storage   deployment.apps/local-path-provisioner   1/1     1            1           3m34s

NAMESPACE            NAME                                                DESIRED   CURRENT   READY   AGE
kube-system          replicaset.apps/coredns-558bd4d5db                  2         2         2       3m23s
local-path-storage   replicaset.apps/local-path-provisioner-547f784dff   1         1         1       3m23s
```

### Install Version One of the Application

Everything looks good. Let's deploy v1 of the app. Apply the yaml in the `infra` directory to deploy the applications.


```sh
kubectl apply -f infra/

deployment.apps/mongo created
service/mongo created
configmap/wizards-client-configmap created
deployment.apps/wizards-client created
configmap/wizards-server-configmap created
deployment.apps/wizards-server created
service/wizards-server created
```

Check the pods.

```sh
kubectl get pods

NAME                              READY   STATUS              RESTARTS   AGE
mongo-6844bdfcdd-f9ts4            0/1     ContainerCreating   0          23s
wizards-client-66cfb96b7f-9g5z5   1/1     Running             0          23s
wizards-server-c445d54d6-zx9d9    0/1     Error               2          23s
```

The `wizards-server` pod here is failing because the MongoDB pod is not up. Wait a few minutes for the MongoDB pod to come up and the server pod should recover.

```sh
kubectl get pods 

NAME                              READY   STATUS    RESTARTS   AGE
mongo-6844bdfcdd-tk52m            1/1     Running   0          83s
wizards-client-66cfb96b7f-9g5z5   1/1     Running   0          83s
wizards-server-c445d54d6-2phq2    1/1     Running   1          83s
```

Version one of the application has been deployed.

### Check the Running Application

The only way for an external user to interact with the application is to tail the logs. Let's look at the server logs first by running `kubectl logs <pod-name>.

```sh
kubectl logs wizards-server-c445d54d6-2phq2

...
2021/06/15 18:25:11 # -------------------------------------- #
2021/06/15 18:25:11 sending list of wizards to client
2021/06/15 18:25:11 # -------------------------------------- #
2021/06/15 18:25:11 sending wizard to client: Harry Potter
2021/06/15 18:25:11 sending wizard to client: Ron Weasley
2021/06/15 18:25:11 sending wizard to client: Hermione Granger
2021/06/15 18:25:11 sending wizard to client: Cho Chang
2021/06/15 18:25:11 sending wizard to client: Luna Lovegood
2021/06/15 18:25:11 sending wizard to client: Sybill Trelawney
2021/06/15 18:25:11 sending wizard to client: Pomona Sprout
2021/06/15 18:25:11 sending wizard to client: Cedric Diggory
2021/06/15 18:25:11 sending wizard to client: Newton Scamander
2021/06/15 18:25:11 sending wizard to client: Cho Chang
2021/06/15 18:25:11 sending wizard to client: Draco Malfoy
2021/06/15 18:25:11 sending wizard to client: Bellatrix Lestrange
2021/06/15 18:25:11 sending wizard to client: Severus Snape
2021/06/15 18:25:11 # -------------------------------------- #
2021/06/15 18:25:11 done sending list of wizards to client
2021/06/15 18:25:11 # -------------------------------------- #
```

Here is a complete batch of wizards retrieved from the database and sent to the client upon request. Let's look at the client by accessing its logs.

```sh
kubectl logs wizards-client-66cfb96b7f-9g5z5

2021/06/15 18:29:24 # -------------------------------------- #
2021/06/15 18:29:24 requesting list of wizards from the server
2021/06/15 18:29:24 # -------------------------------------- #
2021/06/15 18:29:24 wizard received: Harry Potter
2021/06/15 18:29:24 wizard received: Ron Weasley
2021/06/15 18:29:24 wizard received: Hermione Granger
2021/06/15 18:29:24 wizard received: Cho Chang
2021/06/15 18:29:24 wizard received: Luna Lovegood
2021/06/15 18:29:24 wizard received: Sybill Trelawney
2021/06/15 18:29:24 wizard received: Pomona Sprout
2021/06/15 18:29:24 wizard received: Cedric Diggory
2021/06/15 18:29:24 wizard received: Newton Scamander
2021/06/15 18:29:24 wizard received: Cho Chang
2021/06/15 18:29:24 wizard received: Draco Malfoy
2021/06/15 18:29:24 wizard received: Bellatrix Lestrange
2021/06/15 18:29:24 wizard received: Severus Snape
2021/06/15 18:29:24 # --------------------------------------#
2021/06/15 18:29:24 received list of wizards from the server
2021/06/15 18:29:24 # ------------------------------------- #
```

Here is a complete batch of logging from the response sent back from the server. This process happens every second. 

The application is in good health. Now let's make some modifications.

### Accessing the Application with Telepresence

This application by default is not accessable from outside of the cluster. Looking at the services, they are all of type [ClusterIP](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) which means the pods behind those services are only accessable within the cluster. 


```sh
kubectl get svc

NAME             TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
kubernetes       ClusterIP   10.96.0.1      <none>        443/TCP     19m
mongo            ClusterIP   10.96.218.75   <none>        27017/TCP   15m
wizards-server   ClusterIP   10.96.57.15    <none>        9999/TCP    15m
```

Those 10.x IP addresses are only local to the cluster. Trying to access the MongoDB service via its IP or DNS name from outside of the cluster times out.

```sh
curl -m 1 10.96.218.75:27017
curl: (28) Connection timed out after 1001 milliseconds

curl -m 1 mongo.default.svc.cluster.local:27017
curl: (6) Could not resolve host: mongo.default.svc.cluster.local
```
In the case of the DNS name, the local DNS could not resolve the requested DNS name. That name only lives within cluster DNS. Running either the server or the client application will fail unless they are both running at the sametime outside of the cluster and if there is an installation of MongoDB outside of the cluster. However, the problems there are:

* all dependencies will have to be installed on the developer's workstation
* dependencies change and there could be environment drift from local development to a deployed cluster
* there may be other cluster resources not available on the workstation
* a local database may have outdated or wrong data

These are to name a few. 

Telepresence alleviates all the above problems. and sets up the networking to make it appear that the laptop is in the cluster and has access to the same dependencies as the deployed applications do. 

Connect to the cluster with Telepresence. 

```sh
telepresence connect

Launching Telepresence Daemon v2.3.0 (api v3)
Connecting to traffic manager...
Connected to context kind-kind (https://127.0.0.1:36891)
```

This command creates a new [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) named `ambassador` and installs everything Telepresence needs on the cluster side.

```sh
kubectl get all -n ambassador

NAME                                   READY   STATUS    RESTARTS   AGE
pod/traffic-manager-759568df76-vbx27   1/1     Running   0          94s

NAME                      TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
service/traffic-manager   ClusterIP   None         <none>        8081/TCP   94s

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/traffic-manager   1/1     1            1           94s

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/traffic-manager-759568df76   1         1         1       94s
```

Now, let's try to connect to MongoDB.

```sh
curl -m 1 mongo.default.svc.cluster.local:27017
It looks like you are trying to access MongoDB over HTTP on the native driver port.
```

We now get a response from MongoDB and can use its DNS name. It is as if the local workstation is in the cluster. We can now make local modifications to the application code and validate them against a deployed environement. 

### Intercept the Server-Bound Traffic and Run the Server Locally

Let's first make an update to the server. The change made will be updating the log messages on the `List` GRPC service.

First, install the Go dependencies for the server.

Install the protobuf compiler from these [instructions](https://grpc.io/docs/protoc-installation/) for your OS/Arch. 

Install Go specific protobuf plugins.

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Install the applications' dependencies.

```sh
go mod download 
```

Finally, build the protobuf files using the included Make target.

```sh
make proto
```

Start the server.

```sh
go run cmd/server/main.go

2021/06/15 14:50:51 starting wizards server
2021/06/15 14:50:51 checking for configMap file at: ./wizards-server-configMap.txt
2021/06/15 14:50:51 did not find config file: ./wizards-server-configMap.txt. using Kubernetes environment
2021/06/15 14:50:51 printing MONGO_HOST: 
2021/06/15 14:50:51 dropping the wizards collection and seeding database
2021/06/15 14:50:52 error pinging DB: server selection error: context deadline exceeded, current topology: { Type: Unknown, Servers: [{ Addr: :27017, Type: Unknown, Average RTT: 0, Last error: connection() error occured during connection handshake: dial tcp :27017: connect: connection refused }, ] }
exit status 1
```

The server immediately starts up and fails. Logs show the server doesn't know where the database is. Let's intercept cluster traffic bound for the deployed server and create a config file based on the environment of the deployed server so that the local config matches exactly with the deployed configuration. This is done with Telepresence.

```sh
telepresence intercept wizards-server --namespace default --port 9999:9999 --env-file wizards-server-configMap.txt

An update of telepresence from version 2.3.0 to 2.3.1 is available. Please visit https://www.getambassador.io/docs/telepresence/latest/install/upgrade/ for more info.
Using Deployment wizards-server
intercepted
    Intercept name         : wizards-server-default
    State                  : ACTIVE
    Workload kind          : Deployment
    Destination            : 127.0.0.1:9999
    Service Port Identifier: 9999
    Volume Mount Error     : sshfs is not installed on your local machine
    Intercepting           : all TCP connections
```

The command tells Kubernetes to intercept any traffic for the `wizards-server` deployment in the `default` namespace at port `9999`, proxy it to the local workstation on port `9999`, and to take the environment of the pods in the `wizards-server` deployment and make a local file [wizards-server-configMap.txt](./wizards-server-configMap.txt). Open this file and notice there are a bunch of key=value pairs. These are the environment variables inside of the `wizard-server` pod. We will use this file later.

Now, requests from the client are going to the local workstation and are going unanswered. Tail the logs from the client. 

```sh
kubectl logs -f wizards-client-66cfb96b7f-9g5z5 

2021/06/15 18:54:47 # --------------------------------------#
2021/06/15 18:54:47 received list of wizards from the server
2021/06/15 18:54:47 # ------------------------------------- #
2021/06/15 18:54:48 # -------------------------------------- #
2021/06/15 18:54:48 requesting list of wizards from the server
2021/06/15 18:54:48 # -------------------------------------- #
2021/06/15 18:54:49 wizard received: Harry Potter
2021/06/15 18:54:49 wizard received: Ron Weasley
2021/06/15 18:54:49 wizard received: Hermione Granger
2021/06/15 18:54:49 wizard received: Cho Chang
2021/06/15 18:54:49 wizard received: Luna Lovegood
2021/06/15 18:54:49 wizard received: Sybill Trelawney
2021/06/15 18:54:49 wizard received: Pomona Sprout
2021/06/15 18:54:49 wizard received: Cedric Diggory
2021/06/15 18:54:49 wizard received: Newton Scamander
2021/06/15 18:54:49 wizard received: Cho Chang
2021/06/15 18:54:49 wizard received: Draco Malfoy
2021/06/15 18:54:49 wizard received: Bellatrix Lestrange
2021/06/15 18:54:49 wizard received: Severus Snape
2021/06/15 18:54:49 # --------------------------------------#
2021/06/15 18:54:49 received list of wizards from the server
2021/06/15 18:54:49 # ------------------------------------- #
2021/06/15 18:54:50 # -------------------------------------- #
2021/06/15 18:54:50 requesting list of wizards from the server
2021/06/15 18:54:50 # -------------------------------------- #

```

The logs have stopped updating every second because the client reaches out on port `9999`, but nothing on the local workstation is listening on port `9999`. Open up another terminal and start the server.

```sh
go run cmd/server/main.go

2021/06/15 15:00:15 starting wizards server
2021/06/15 15:00:15 checking for configMap file at: ./wizards-server-configMap.txt
2021/06/15 15:00:15 found ./wizards-server-configMap.txt file. setting environment variables
2021/06/15 15:00:15 done setting environment variables
2021/06/15 15:00:15 printing MONGO_HOST: mongo.default.svc.cluster.local
2021/06/15 15:00:15 dropping the wizards collection and seeding database
2021/06/15 15:00:15 connected to the database
2021/06/15 15:00:15 inserted document for wizard: Harry Potter
2021/06/15 15:00:15 inserted document for wizard: Ron Weasley
2021/06/15 15:00:15 inserted document for wizard: Hermione Granger
2021/06/15 15:00:15 inserted document for wizard: Cho Chang
2021/06/15 15:00:15 inserted document for wizard: Luna Lovegood
2021/06/15 15:00:15 inserted document for wizard: Sybill Trelawney
2021/06/15 15:00:15 inserted document for wizard: Pomona Sprout
2021/06/15 15:00:15 inserted document for wizard: Cedric Diggory
2021/06/15 15:00:15 inserted document for wizard: Newton Scamander
2021/06/15 15:00:15 inserted document for wizard: Cho Chang
2021/06/15 15:00:15 inserted document for wizard: Draco Malfoy
2021/06/15 15:00:15 inserted document for wizard: Bellatrix Lestrange
2021/06/15 15:00:15 inserted document for wizard: Severus Snape
2021/06/15 15:00:15 inserted 13 documents
2021/06/15 15:00:15 finished seeding the database
2021/06/15 15:00:15 # ----------------------------------- #
2021/06/15 15:00:15 running grpc server on port: 9999
2021/06/15 15:00:15 # ----------------------------------- #
```

Now the server successfully starts up because upon start up, it read the environment file and got the same configs that are available to the deployed application. 

Notice how after a few seconds the client started updating every second and the server logs started updating? The deployed client is now able to access the locally running Go process. If you were to open another terminal and tail the logs of the deployed server, they would not be updating. It receives no traffic, while the local process does. 

Great, now let's update the local process.

### Update the Server

There needs to be an update to the server. These logs are fine, but users complain that the logs need more flair. Let's respond to this feature request. 

Open the [server](./cmd/server/main.go) file. There's a for loop at lines `110-137` with commented out code:

```go
	for cursor.Next(context.Background()) {
		var wizard pb.Wizard

		err := cursor.Decode(&wizard)
		if err != nil {
			log.Printf("unable to decode wizard cursor into struct: %v", err)
		}

		// Commenting out. Uncomment when making the server change in the demo
		// switch house := wizard.House; house {
		// case "Gryffindor":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Eagle, wizard.GetName())
		// case "Ravenclaw":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Bird, wizard.GetName())
		// case "Hufflepuff":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Badger, wizard.GetName())
		// case "Slytherin":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Snake, wizard.GetName())
		// }

		// comment this log statement as part of the server demo
		log.Printf("sending wizard to client: %v", wizard.GetName())

		err = srv.Send(&wizard)
		if err != nil {
			log.Printf("send error: %v", err)
		}
	}
```

Let's uncomment the commented lines and comment out line 131. It should look like this:

```go
	for cursor.Next(context.Background()) {
		var wizard pb.Wizard

		err := cursor.Decode(&wizard)
		if err != nil {
			log.Printf("unable to decode wizard cursor into struct: %v", err)
		}

		// Commenting out. Uncomment when making the server change in the demo
		switch house := wizard.House; house {
		case "Gryffindor":
			log.Printf("%v - sending wizard to client: %v", emoji.Eagle, wizard.GetName())
		case "Ravenclaw":
			log.Printf("%v - sending wizard to client: %v", emoji.Bird, wizard.GetName())
		case "Hufflepuff":
			log.Printf("%v - sending wizard to client: %v", emoji.Badger, wizard.GetName())
		case "Slytherin":
			log.Printf("%v - sending wizard to client: %v", emoji.Snake, wizard.GetName())
		}

		// comment this log statement as part of the server demo
		// log.Printf("sending wizard to client: %v", wizard.GetName())

		err = srv.Send(&wizard)
		if err != nil {
			log.Printf("send error: %v", err)
		}
	}

```

You will need to install and import the `emoji` package: `go get github.com/enescakir/emoji && go mod tidy`. 

Restart the server by pressing `ctrl+c` and running the `go run` command again.

Now the logs print out some ASCII emojis based on the wizard's house they belong to:

```sh
2021/06/15 15:09:43 # -------------------------------------- #
2021/06/15 15:09:43 sending list of wizards to client
2021/06/15 15:09:43 # -------------------------------------- #
2021/06/15 15:09:43 🦅 - sending wizard to client: Harry Potter
2021/06/15 15:09:43 🦅 - sending wizard to client: Ron Weasley
2021/06/15 15:09:43 🦅 - sending wizard to client: Hermione Granger
2021/06/15 15:09:43 🐦 - sending wizard to client: Cho Chang
2021/06/15 15:09:43 🐦 - sending wizard to client: Luna Lovegood
2021/06/15 15:09:43 🐦 - sending wizard to client: Sybill Trelawney
2021/06/15 15:09:43 🦡 - sending wizard to client: Pomona Sprout
2021/06/15 15:09:43 🦡 - sending wizard to client: Cedric Diggory
2021/06/15 15:09:43 🦡 - sending wizard to client: Newton Scamander
2021/06/15 15:09:43 🐍 - sending wizard to client: Draco Malfoy
2021/06/15 15:09:43 🐍 - sending wizard to client: Bellatrix Lestrange
2021/06/15 15:09:43 🐍 - sending wizard to client: Severus Snape
2021/06/15 15:09:43 # -------------------------------------- #
2021/06/15 15:09:43 done sending list of wizards to client
2021/06/15 15:09:43 # -------------------------------------- #
```

This is a much better user experience and fulfills the feature reqeust. If you still have the deployed server terminal open, it should not be updating, nor should the client logs be impacted. 

Telepresence intercept is a great way to make local changes to an application and validate them with deployed applications, configurations, and data. It is far more likely that a regression could be caught in development since there is the ability to do some integration testing right on the developer's workstation before the code is even deployed. 

Let's clean up.

Stop the server process with `ctrl+c`. 

List out the Telepresence intercepts.

```sh
telepresence list

mongo         : ready to intercept (traffic-agent not yet installed)
wizards-server: intercepted
    Intercept name         : wizards-server-default
    State                  : ACTIVE
    Workload kind          : Deployment
    Destination            : 127.0.0.1:9999
    Service Port Identifier: 9999
    Intercepting           : all TCP connections
```

This command shows all applications that can be intercepted. In this case, it is all of the deployments who have a service associated with them. Since the client does not have a service, it cannot (nor does it need to be) intercepted.

Disconnect from the server and list the connections.

```sh
telepresence leave wizards-server-default

telepresence list

mongo         : ready to intercept (traffic-agent not yet installed)
wizards-server: ready to intercept (traffic-agent already installed)
```

Now both available intercepts are ready to be started. You should also notice that both the deployed client and deployed server have begun updating their logs again.

Now, let's update the client and validate against the deployed server.

### Update the Client Locally and Validate Against the Deployed Server

Running the client now that the `wizards-server-configMap.txt` file has been created is much easier than running the server. There are not intercepts that need to be done. Running the client is a simple command.

```sh
go run cmd/client/main.go
2021/06/15 15:21:06 checking for configMap file at: ./wizards-server-configMap.txt
2021/06/15 15:21:06 found ./wizards-server-configMap.txt file. setting environment variables
2021/06/15 15:21:06 done setting environemnt variables
2021/06/15 15:21:07 # -------------------------------------- #
2021/06/15 15:21:07 requesting list of wizards from the server
2021/06/15 15:21:07 # -------------------------------------- #
2021/06/15 15:21:07 wizard received: Harry Potter
2021/06/15 15:21:07 wizard received: Ron Weasley
2021/06/15 15:21:07 wizard received: Hermione Granger
2021/06/15 15:21:07 wizard received: Cho Chang
2021/06/15 15:21:07 wizard received: Luna Lovegood
2021/06/15 15:21:07 wizard received: Sybill Trelawney
2021/06/15 15:21:07 wizard received: Pomona Sprout
2021/06/15 15:21:07 wizard received: Cedric Diggory
2021/06/15 15:21:07 wizard received: Newton Scamander
2021/06/15 15:21:07 wizard received: Cho Chang
2021/06/15 15:21:07 wizard received: Draco Malfoy
2021/06/15 15:21:07 wizard received: Bellatrix Lestrange
2021/06/15 15:21:07 wizard received: Severus Snape
2021/06/15 15:21:07 # --------------------------------------#
2021/06/15 15:21:07 received list of wizards from the server
2021/06/15 15:21:07 # ------------------------------------- #
```

The client starts up, finds the config file, and begins asking the deployed server for a list of wizards. 

Let's make a change to this local client. Our users have asked us to alert them if a wizard in the return list is a Death Eater. We can deliver this reqeust easily since Mongo already has that data stored. 

Open the [client](./cmd/client/main.go) file.

There are a few commented out lines at `75-77` and a whole function is commented out at `82-86`. This will enable to functionality to update the logs if a Death Eater is found in the response from the server.

```go
		// commenting out for the demo. uncomment during demo of the client
		// if resp.GetDeathEater() {
		// 	alertDeathEather(resp)
		// }
	}
}

// commenting out for the demo. uncomment during demo of the client
// func alertDeathEather(w *pb.Wizard) {
// 	log.Println("")
// 	log.Printf("Oh no! %s is a Death Eater!", w.GetName())
// 	log.Println("")
// }

```

Uncomment the lines to look like this:

```go
		// commenting out for the demo. uncomment during demo of the client
		if resp.GetDeathEater() {
			alertDeathEather(resp)
		}
	}
}

// commenting out for the demo. uncomment during demo of the client
func alertDeathEather(w *pb.Wizard) {
	log.Println("")
	log.Printf("Oh no! %s is a Death Eater!", w.GetName())
	log.Println("")
}
```

Save the file.

Restart the client by pressing `ctrl+c` and running `go run cmd/client/main.go`

```sh
2021/06/15 15:27:30 checking for configMap file at: ./wizards-server-configMap.txt
2021/06/15 15:27:30 found ./wizards-server-configMap.txt file. setting environment variables
2021/06/15 15:27:30 done setting environemnt variables
2021/06/15 15:27:31 # -------------------------------------- #
2021/06/15 15:27:31 requesting list of wizards from the server
2021/06/15 15:27:31 # -------------------------------------- #
2021/06/15 15:27:31 wizard received: Harry Potter
2021/06/15 15:27:31 wizard received: Ron Weasley
2021/06/15 15:27:31 wizard received: Hermione Granger
2021/06/15 15:27:31 wizard received: Cho Chang
2021/06/15 15:27:31 wizard received: Luna Lovegood
2021/06/15 15:27:31 wizard received: Sybill Trelawney
2021/06/15 15:27:31 wizard received: Pomona Sprout
2021/06/15 15:27:31 wizard received: Cedric Diggory
2021/06/15 15:27:31 wizard received: Newton Scamander
2021/06/15 15:27:31 wizard received: Cho Chang
2021/06/15 15:27:31 wizard received: Draco Malfoy
2021/06/15 15:27:31 
2021/06/15 15:27:31 Oh no! Draco Malfoy is a Death Eater!
2021/06/15 15:27:31 
2021/06/15 15:27:31 wizard received: Bellatrix Lestrange
2021/06/15 15:27:31 
2021/06/15 15:27:31 Oh no! Bellatrix Lestrange is a Death Eater!
2021/06/15 15:27:31 
2021/06/15 15:27:31 wizard received: Severus Snape
2021/06/15 15:27:31 # --------------------------------------#
2021/06/15 15:27:31 received list of wizards from the server
2021/06/15 15:27:31 # ------------------------------------- #

```
Now the client starts up, requests a list of wizards from the server, and if it finds that a wizard is a Death Eater, there is a warning printed in the logs. 

Good. Our users will be happy with this change.

This concludes the demo. 

### Clean Up

We should be good stewards and clean up. 

First, disconnect Telepresence from the cluster.

```sh
telepresence quit

Telepresence Daemon quitting...done

telepresence status

Root Daemon: Not running
User Daemon: Not running
```

Telepresence is now no longer connected the workstation to the cluster.

The cluster resources can now be cleaned up.

```sh
kind delete cluster

Deleting cluster "kind" ...
```

This concludes the demo.

## Contributing

This repository is not accepting contrubutions. 