# Deploy Applications and Services to GKE clusters

Repeat following steps for each of your clusters.

1. Deploy Application

```console
# first cluster
kubectl apply -f deploy-1.0.yaml
# second cluster
kubectl apply -f deploy-2.0.yaml
```

2. Create K8s Services for Application

```console
# first cluster
kubectl apply -f svc-1.0.yaml
# second cluster
kubectl apply -f svc-2.0.yaml
```

Note the `cloud.google.com/neg: '{"exposed_ports": {"80":{}}}'` annotation on the service telling GKE to create a NEG for the service.

# Setup Load Balancing (GCLB) Components

1. Create Health Check:

```console
gcloud compute health-checks create http hello-world-healthz \
  --use-serving-port \
  --request-path="/"
```

2. Create Backend Services

Create backend service for each of the services, plus one more to serve as default backend for traffic that doesn’t match the path-based rules.

```console
gcloud compute backend-services create backend-service-default \
  --global

gcloud compute backend-services create backend-service-hellow-world \
  --global \
  --health-checks hello-world-healthz
```

3. Create URL Map (Load Balancer)

```console
gcloud compute url-maps create unified-load-balancer \
  --global \
  --default-service backend-service-default
```

4. Add Path Rules to URL Map

```console
gcloud compute url-maps add-path-matcher unified-load-balancer \
  --global \
  --path-matcher-name=hello-world-matcher \
  --default-service=backend-service-default \
  --backend-service-path-rules='/hello-world/*=backend-service-hellow-world'
```

5. Reserve Static IP Address

```console
gcloud compute addresses create unified-load-balancer-ipv4 \
  --ip-version=IPV4 \
  --global
```

6. Create Target HTTP Proxy

```console
gcloud compute target-http-proxies create unified-http-proxy \
  --url-map=unified-load-balancer
```

7. Create Forwarding Rule

```console
gcloud compute forwarding-rules create hellow-world-fw-rule \
  --target-http-proxy=unified-http-proxy \
  --global \
  --ports=80 \
  --address=unified-load-balancer-ipv4
```

# Connect K8s Services to the Load Balancer

GKE has provisioned NEGs for each of the K8s services deployed with the `cloud.google.com/neg` annotation. Now we need to add these NEGs as backends to corresponding backend services.

1. Retrieve Names of Provisioned NEGs

```console
kubectl get svc \
  -o custom-columns='NAME:.metadata.name,NEG:.metadata.annotations.cloud\.google\.com/neg-status'
```

Note down the NEG name and zones for each service.

Repeat for all your GKE Clusters.

2. Add NEGs to Backend Services

Repeat following for every NEG and zone from both clusters. Make sure to use only NEGs belonging to the same service.

```console
gcloud compute backend-services add-backend backend-service-hellow-world \
 --global \
 --network-endpoint-group [neg_name] \
 --network-endpoint-group-zone=[neg_zone] \
 --balancing-mode=RATE \
 --max-rate-per-endpoint=100
```

3. Allow GCLB Traffic

```console
gcloud compute firewall-rules create fw-allow-unified-load-balancer \
  --network=[vpc_name] \
  --action=allow \
  --direction=ingress \
  --source-ranges=130.211.0.0/22,35.191.0.0/16 \
  --rules=tcp:8080
```

# Test Everything’s Working

```console
# replace 34.117.111.28 with your static IP
curl -v http://34.117.111.28/hello-world/
```

# References

* [Multi-Cluster Load Balancing with GKE](https://stepan.wtf/multi-cluster-load-balancing-with-gke/)
