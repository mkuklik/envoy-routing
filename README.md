# envoy-routing

1. envoy with routing table splitting traffic between two servers

2. envoy with routing table pulled from xDS splitting traffic between two servers

# minikube 
```bash
minikube start
minikube stop
minikube dashboard
```
## enable ingress
```
minikube addons enable ingress
```

## swagger for api server
```bash
minikube start --extra-config=apiserver.Features.EnableSwaggerUI=true
# proxy
kubectl proxy --port=8080
# web  http://localhost:8080/swagger-ui/

```

## start proxy 
```
kubectl proxy --port=8080
```

To access api: 
```bash
curl http://localhost:8080/api/
```

# skaffold
skaffold dev
skaffold run


helm repo add slamdev-helm-charts https://slamdev.github.io/helm-charts
helm install slamdev-helm-charts/envoy --generate-name


# Envoy on minikube
https://www.getambassador.io/resources/envoy-flask-kubernetes/
