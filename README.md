# Telepresence Demo

This repository holds code for a demonstration of [Telepresence](https://www.telepresence.io/). Below describes the scripts and commands to setup the demo and run through it. 

The demonstration has three components, all running inside of a Kubernetes cluster, that are not accessable from outside of the cluster (which is the point of Telepresence):

* [MongoDB](https://docs.mongodb.com/)
* [GRPC Client](./cmd/client/main.go)
* [GRPC Server](./cmd/server/main.go)


## The Applications

The demo applications are simple. They are [GRPC](https://grpc.io/docs/) apps are written in [Go](https://golang.org).


The application is a wizards listing service as defined in the [wizards.proto](./wizards/wizards.proto) file. The client requests a list of wizards from the server. The server retrieves a list of wizards from the database and [streams](https://grpc.io/docs/what-is-grpc/core-concepts/#server-streaming-rpc) them back to the client.

*Note: The apps are meant to be a quickly written simple single file. There are obviously a lot improvements to be made and production best practices implimented. They are just single files with a main package and they import protobuf. The goal of this repo is not to write the prettiest or best tested Go programs, but to demo Telepresence and how it can improve development workflows.

## Non-Application Dependencies

* [Telepresence CLI](https://www.telepresence.io/docs/latest/install/)
* [Kubernetes](https://kubernetes.io/) (using [Kind](https://kind.sigs.k8s.io/))
* [GNI Make](https://www.gnu.org/software/make/)

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
wizards-client-66cfb96b7f-bd4w5   1/1     Running             0          23s
wizards-server-c445d54d6-zx9d9    0/1     Error               2          23s
```

The `wizards-server` pod here is failing becuase the MongoDB pod is not up. Wait a few minutes for the MongoDB pod to come up and the server pod should recover.

```sh
kubectl get pods 

NAME                              READY   STATUS    RESTARTS   AGE
mongo-6844bdfcdd-tk52m            1/1     Running   0          83s
wizards-client-66cfb96b7f-v4fxp   1/1     Running   0          83s
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

Here is a complete batch of wizards retreived from the database and sent to the cleint upon reqeust. Let's look at the client by accessing its logs.

```sh

```