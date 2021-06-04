# Telepresence Demo

## Protobuff

Build the protobuff definition

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto
```

## Working with MongoDB

### Deploy MongoDB

Apply the yamls

```sh
kubectl apply -f infra/
```

After a few minutes, the Mongo pod should be up
```sh
NAME                         READY   STATUS    RESTARTS   AGE
pod/mongo-7c5d9d7659-cxkxw   1/1     Running   0          143m

NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
service/kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP     4h31m
service/mongo        ClusterIP   10.96.47.222   <none>        27017/TCP   142m

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mongo   1/1     1            1           143m

NAME                               DESIRED   CURRENT   READY   AGE
replicaset.apps/mongo-7c5d9d7659   1         1         1       143m

```
### Connect to MongoDB

1. Connect to the running pod
```sh
kubectl exec -ti <pod name> -- bash
```

2. Connect to MongoDB using the builtin CLI
    ```sh
    mongo
    ```