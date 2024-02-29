# Image URL to use all building/pushing image targets
IMG ?= controller:latest

INVENTORY_BIN_NAME = "inventory"
LLDP_UPDATE_BIN_NAME = "nic-updater"
BENCHMARK_BIN_NAME = "benchmark"
BENCHMARK_SCHEDULER_BIN_NAME = "benchmark-scheduler"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec
.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: goimports
	go fmt ./...
	$(GOIMPORTS) -w .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet envtest ## Run tests.
	@echo "--> Project testing"
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test -race ./... -coverprofile cover.out
	@echo "--> Making coverage html page"
	go tool cover -html cover.out -o ./index.html

.PHONY: add-license
add-license: addlicense ## Add license header to all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -f hack/license-header.txt {} +

.PHONY: check-license
check-license: addlicense ## Check license header presence in all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -check -c 'IronCore authors' {} +

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint
	$(GOLANGCI_LINT) run --fix

.PHONY: check
check: vet add-license lint test # Run lint, add licenses, test

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	for BIN_NAME in $(INVENTORY_BIN_NAME) $(LLDP_UPDATE_BIN_NAME) $(BENCHMARK_BIN_NAME) $(BENCHMARK_SCHEDULER_BIN_NAME); do \
		go build -o dist/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
	done
	cp -rf res/ dist/

.PHONY: dl-pciids
dl-pciids:
	curl https://pci-ids.ucw.cz/v2.2/pci.ids --output ./res/pci.ids

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./main.go

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

# PLATFORMS defines the target platforms for  the manager image be build to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - able to use docker buildx . More info: https://docs.docker.com/build/buildx/
# - have enable BuildKit, More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image for your registry (i.e. if you do not inform a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To properly provided solutions that supports more than one platform you should use this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: test ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- docker buildx create --name project-v3-builder
	docker buildx use project-v3-builder
	- docker buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- docker buildx rm project-v3-builder
	rm Dockerfile.cross

.PHONY: clean
clean:
	rm -rf ./dist/

##@ Build Dependencies

## Location to install dependencies to
LOCAL_BIN ?= $(shell pwd)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

## Tool Binaries
ADDLICENSE ?= $(LOCAL_BIN)/addlicense
GOLANGCI_LINT ?= $(LOCAL_BIN)/golangci-lint
GOIMPORTS ?= $(LOCAL_BIN)/goimports
ENVTEST ?= $(LOCAL_BIN)/setup-envtest
KUSTOMIZE ?= $(LOCAL_BIN)/kustomize

## Tool Versions
ADDLICENSE_VERSION ?= v1.1.1
GOLANGCI_LINT_VERSION ?= v1.55.2
GOIMPORTS_VERSION ?= v0.16.1
ENVTEST_K8S_VERSION ?= 1.28.3
KUSTOMIZE_VERSION ?= v4.5.4

.PHONY: addlicense
addlicense: $(ADDLICENSE)
$(ADDLICENSE): $(LOCAL_BIN)
	@test -s $(ADDLICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCAL_BIN)
	@test -s $(GOLANGCI_LINT) && $(GOLANGCI_LINT) --version | grep -q $(GOLANGCI_LINT_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: goimports
goimports: $(GOIMPORTS)
$(GOIMPORTS): $(LOCAL_BIN)
	@test -s $(GOIMPORTS) || GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCAL_BIN)
	@test -s $(ENVTEST) || GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCAL_BIN)
	@test -s $(KUSTOMIZE) || GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/kustomize/kustomize/v4@$(KUSTOMIZE_VERSION)
