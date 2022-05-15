[![Go Report Card](https://goreportcard.com/badge/github.com/akankshakumari393/depkon)](https://goreportcard.com/report/github.com/akankshakumari393/depkon)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
# Depkon Kubernetes Operator

This operator tries to sync a configmap resorces with a list of deployment for a particular namespace. In this project we have defined a Custom Resource `depkon`.

# Install using helm

```
helm repo add akankshakumari393 https://akankshakumari393.github.io/helm-charts
kubectl create namespace controller
helm install depkon akankshakumari393/depkon-operator -n controller
```
