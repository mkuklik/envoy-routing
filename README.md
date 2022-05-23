# envoy-routing

1. envoy with routing table splitting traffic between two servers

2. envoy with routing table pulled from xDS splitting traffic between two servers

3. envoy internal proxy with routing table pulled from xDS splitting traffic between two servers

4. envoy two stage internal proxy

5. envoy multicasting

# minikube 
```bash
minikube start
minikube stop
minikube dashboard
```
## access service directly
```
minikube service backend-service
```
## access individual pod
```
kubectl port-forward <POD NAME> 5000
kubectl port-forward <POD NAME> 8080:5000

kubectl port-forward backend-56b857cb7-cvlxb 5000
kubectl port-forward backend-56b857cb7-cvlxb 8080:5000
```
## use loadbalance service
```
minikube service lb-srv
```

## Ingress
### enable ingress
```
minikube addons enable ingress
minikube ip
```
### ingress ip address
```
minikube ip
curl $(minikube ip)
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

## skaffold
skaffold dev
skaffold run
kubectl config current-context

skaffold run --tail -f <file> -p <profile>

helm repo add slamdev-helm-charts https://slamdev.github.io/helm-charts
helm install slamdev-helm-charts/envoy --generate-name


# Envoy on minikube
https://www.getambassador.io/resources/envoy-flask-kubernetes/
https://github.com/datawire/envoy-steps/




Tasks:
- run this with regular ingress and check backend working
- add envoy
- add static routing 
- add xDS dynamic routing
