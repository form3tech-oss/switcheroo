To test locally (minikube (or kind)) use cert manager to manage webhook certificates.

1. To install cert manager: https://cert-manager.io/docs/installation/kubernetes/
with regular manifest, use: kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.3.1/cert-manager.yaml

2. After, apply the `issuer.yaml` to have a generic self-signed cert issuer on the cluster.

3. Deploy the application with helm chart.
   1. use `eval $(minikube docker-env)` for minikube test
   2. build the image with `make docker-build`
   3. Set also `imagePullPolicy: Never` in template

