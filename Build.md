### Using code-generator to generate listers, informers, clientset

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
.
├── Build.md
├── controller
│   └── depkon.go
├── Dockerfile
├── go.mod
├── go.sum
├── hack
│   ├── boilerplate.go.txt
│   ├── tools.go
│   ├── update-codegen.sh
│   └── verify-codegen.sh
├── helm
│   └── depkon-operator
│       ├── charts
│       ├── Chart.yaml
│       ├── crds
│       │   └── akankshakumari393.dev_depkons.yaml
│       ├── templates
│       │   ├── deployment.yaml
│       │   ├── _helpers.tpl
│       │   ├── rbac.yaml
│       │   └── serviceaccount.yaml
│       └── values.yaml
├── main.go
├── manifests
│   ├── akankshakumari393.dev_depkonlists.yaml
│   ├── akankshakumari393.dev_depkons.yaml
│   └── depkon-cr.yaml
├── pkg
│   ├── apis
│   │   └── akankshakumari393.dev
│   │       └── v1alpha1
│   │           ├── doc.go
│   │           ├── register.go
│   │           ├── types.go
│   │           └── zz_generated.deepcopy.go
│   └── generated
│       ├── clientset
│       │   └── versioned
│       │       ├── clientset.go
│       │       ├── doc.go
│       │       ├── fake
│       │       ├── scheme
│       │       └── typed
│       ├── informers
│       │   └── externalversions
│       │       ├── akankshakumari393.dev
│       │       ├── factory.go
│       │       ├── generic.go
│       │       └── internalinterfaces
│       └── listers
│           └── akankshakumari393.dev
│               └── v1alpha1
├── README.md
└── vendor
    ├── github.com
    ├── golang.org
    ├── google.golang.org
    ├── gopkg.in
    ├── k8s.io
    ├── modules.txt
    └── sigs.k8s.io
```
### Using controller-gen to generate CRD specifications/yaml manifests

```
{relative_path_to_controller_tools_directory}/controller-tools/cmd/controller-gen/controller-gen paths=github.com/akankshakumari393/depkon/pkg/apis/akankshakumari393.dev/v1alpha1 crd:crdVersions=v1 output:crd:artifacts:config=manifests
```

Note: group name should be a domain with at least one dot

### Execute
```
# Build the image 
docker build -t akankshakumari393/depkon:0.0.1 .

# Push the image 
docker push akankshakumari393/depkon:0.0.1

# create a namespace
kubectl create namespace controller

# Install helm chart locally
helm install depkonlocal ./helm/depkon-operator/ -n controller

# Create a Depkon CR in any namespace, this will sync the configmap and deployment in that namespace
kubectl create -f ./manifests/depkon-cr.yaml
```
