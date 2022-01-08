
context:
	kubectl config current-context


sdev:
	SKAFFOLD_DEFAULT_REPO=gcr.io/gaggle-996 skaffold dev

minikube-context:
	minikube start -p envoy-profile
	source <(minikube docker-env -p envoy-profile)
	skaffold config set --kube-context envoy-profile local-cluster true
