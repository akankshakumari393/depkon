# Depkon Kubernetes Operator

## Description

This operator tries to sync a configmap resorces with a list of deployment for a particular namespace. In this project we have defined a Custom Resource `depkon`.

### Using code-generator to generator scaffolding 

Script for code generator is in `hack/`

This must be set if go version > 1.13 

`GO111MODULE="" or GO111MODULE=on`

```
# generate vendor directory containing code-generator, -e is specified to ignore error 
go mod vendor -e
# Run code-generator 
bash ./hack/update-codegen.sh
```
Generated Project structure
```
├── github.com
│   └── akankshakumari393
│       └── depkon
│           ├── depkon
│           ├── go.mod
│           ├── go.sum
│           ├── hack
│           │   ├── boilerplate.go.txt
│           │   ├── tools.go
│           │   ├── update-codegen.sh
│           │   └── verify-codegen.sh
│           ├── main.go
│           ├── manifests
│           │   ├── akankshakumari393.dev_depkonlists.yaml
│           │   └── akankshakumari393.dev_depkons.yaml
│           ├── pkg
│           │   ├── apis
│           │   └── generated
│           ├── README.md
│           └── vendor
│               ├── github.com
│               ├── golang.org
│               ├── google.golang.org
│               ├── gopkg.in
│               ├── k8s.io
│               ├── modules.txt
│               └── sigs.k8s.io
```
### Using controller-gen to generate CRD specifications/yaml manifests

```
{relative_path_to_controller_tools_directory}/controller-tools/cmd/controller-gen/controller-gen paths=github.com/akankshakumari393/depkon/pkg/apis/akankshakumari393.dev/v1alpha1 crd:crdVersions=v1 output:crd:artifacts:config=manifests
```

Note: group name should be a domain with at least one dot

### Build the program 

`go build`

### create CRD

`
kubectl create -f manifests/akankshakumari393.dev_depkons.yaml
`
### Execute
