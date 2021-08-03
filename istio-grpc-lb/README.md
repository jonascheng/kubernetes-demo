# Problem to gRPC Load Balancing

gRPC users are surprised to find that Kubernetes's default load balancing often doesn't work out of the box with gRPC.

gRPC breaks the standard connection-level load balancing, including what's provided by Kubernetes. This is because gRPC is built on HTTP/2, and HTTP/2 is designed to have a single long-lived TCP connection, across which all requests are multiplexed—meaning multiple requests can be active on the same connection at any point in time. Normally, this is great, as it reduces the overhead of connection management. However, it also means that connection-level balancing isn't very useful. Once the connection is established, there's no more balancing to be done. All requests will get pinned to a single destination pod, as shown below:

![](https://d33wubrfki0l68.cloudfront.net/b05eb6c0d5c4672ed795cc1c44ca476987057436/a6cda/images/blog/grpc-load-balancing-with-linkerd/mono-8d2e53ef-b133-4aa0-9551-7e36a880c553.png)

In this demo you will deploy a gRPC service as well as a gRPC client. The gRPC client send requests infintely to observe response from gRPC service. You will  experience the aforementioned issue and solve it with istio in action.

>
> * [gRPC Load Balancing on Kubernetes without Tears](https://kubernetes.io/blog/2018/11/07/grpc-load-balancing-on-kubernetes-without-tears/)
> * [為什麼我們要用Istio，Native Kubernetes有什麼做不到](https://ithelp.ithome.com.tw/articles/10217407)


# Prerequisites

Before you begin, check the following prerequisites:

1. A cluster running a compatible version of Kubernetes (1.18, 1.19, 1.20, 1.21).
2. [Download the Istio release.](https://istio.io/latest/docs/setup/getting-started/#download)
3. Create a namespace `istio-grpc-demo` for this demo deployment.

```console
$> kubectl create ns istio-grpc-demo
```

# With Istio Disabled

Make sure Istio is disabled

```console
$> kubectl label namespace istio-grpc-demo istio-injection=disable --overwrite=true

namespace/istio-grpc-demo labeled
```

Deploy client and server, and make sure `Pod READY 1/1` (without Istio proxy)

```console
$> kubectl create -f grpc-server.yaml
$> kubectl create -f grpc-client.yaml
$> kubectl get po -n istio-grpc-demo

NAME                           READY   STATUS    RESTARTS   AGE
grpc-client-85c48f8dd4-tjztb   1/1     Running   0          20s
grpc-server-676574776f-2mjhb   1/1     Running   0          26s
grpc-server-676574776f-9qzfg   1/1     Running   0          26s
grpc-server-676574776f-qr92n   1/1     Running   0          26s
```

Observe IP responsed from grpc-server

```console
$> 
```

# With Istio Enabled

# Cleanup

```console
$> kubectl delete ns istio-grpc-demo
```