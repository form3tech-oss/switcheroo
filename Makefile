
GOFLAGS = -mod=vendor
export GOFLAGS
export K8S_VERSION=1.19.2
export CGO_ENABLED=0

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: test

prepare-testenv-binaries:
	curl -sSLo envtest-bins.tar.gz "https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-${K8S_VERSION}-linux-amd64.tar.gz"
	mkdir -p test-bin/
	tar -C . --strip-components=1 -zvxf envtest-bins.tar.gz
	rm envtest-bins.tar.gz

# Run tests
.PHONY: test
test: prepare-testenv-binaries fmt vet
	go test ./...

# Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	go vet ./...

# Generate code
.PHONY: generate
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

.PHONY: build
build:
	go build ./cmd/switcheroo/main.go

.PHONY: docker-build
docker-build: build
	docker build . -t switcheroo:latest

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-beta.4
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
