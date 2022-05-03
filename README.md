# Depkon Kubernetes Operator

This operator tries to sync a configmap resorces with a list of deployment for a particular namespace. In this project we have defined a Custom Resource `depkon`.

# Install using helm

```
helm repo add akankshakumari393 https://akankshhakumari393.github.io/helm-charts
kubectl create namespace controller
helm install depkon akankshakumari393/depkon-operator -n controller
```