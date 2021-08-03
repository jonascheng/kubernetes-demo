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

# Without Istio

Make sure Istio is removed and disabled

```console
$> istioctl x uninstall --purge
$> kubectl label namespace istio-grpc-demo istio-injection=disable --overwrite=true

namespace/istio-grpc-demo labeled
```

Deploy client and server, and make sure `Pod READY 1/1` (without Istio proxy)

```console
$> kubectl create -f grpc-server.yaml
$> kubectl create -f grpc-client.yaml
$> kubectl get po -n istio-grpc-demo

NAME                           READY   STATUS    RESTARTS   AGE
grpc-client-75fb696656-9kxsr   1/1     Running   0          43s
grpc-server-676574776f-7s87d   1/1     Running   0          48s
grpc-server-676574776f-bh758   1/1     Running   0          48s
grpc-server-676574776f-pqn44   1/1     Running   0          48s
```

Observe IP responded from grpc-server

```console
# replace grpc-client-75fb696656-9kxsr with corresponding client pod
$> kubectl logs -f -n istio-grpc-demo grpc-client-75fb696656-9kxsr

2021/08/03 10:02:14 Current pod ip: 10.100.2.27
2021/08/03 10:02:16 Resp: 10.100.1.20
2021/08/03 10:02:17 Resp: 10.100.1.20
2021/08/03 10:02:18 Resp: 10.100.1.20
2021/08/03 10:02:19 Resp: 10.100.1.20
2021/08/03 10:02:20 Resp: 10.100.1.20
```

# With Istio Enabled

Cleanup previous deployment and recreate namespace

```console
$> kubectl delete ns istio-grpc-demo; kubectl create ns istio-grpc-demo
```

Enable Istio

```console
$> kubectl label namespace istio-grpc-demo istio-injection=enabled --overwrite=true

namespace/istio-grpc-demo labeled

$> istioctl install

✔ Istio core installed        
✔ Istiod installed  
✔ Ingress gateways installed                                                             
✔ Installation complete
```

Deploy client and server, and make sure `Pod READY 2/2` (with Istio proxy)

```console
$> kubectl create -f grpc-server.yaml
$> kubectl create -f grpc-client.yaml
$> kubectl get po -n istio-grpc-demo

NAME                           READY   STATUS    RESTARTS   AGE
grpc-client-75fb696656-bbcxb   2/2     Running   0          17s
grpc-server-676574776f-2mkj5   2/2     Running   0          22s
grpc-server-676574776f-d2qtv   2/2     Running   0          22s
grpc-server-676574776f-s2b24   2/2     Running   0          22s
```

Observe IP responded from grpc-server

```console
# replace grpc-client-75fb696656-bbcxb with corresponding client pod
$> kubectl logs -f -n istio-grpc-demo grpc-client-75fb696656-bbcxb grpc-client

2021/08/03 10:02:14 Current pod ip: 10.100.2.27
2021/08/03 10:02:16 Resp: 10.100.1.20
2021/08/03 10:02:17 Resp: 10.100.1.20
2021/08/03 10:02:18 Resp: 10.100.1.20
2021/08/03 10:02:19 Resp: 10.100.1.20
2021/08/03 10:02:20 Resp: 10.100.1.20
```

As the result, without Istio the IP responded from server side was constant, on the flip side the IP was randomly responded with Istio.

# Cleanup

```console
$> istioctl x uninstall --purge
$> kubectl delete ns istio-grpc-demo
```