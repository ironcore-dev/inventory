INVENTORY_BIN_NAME = "inventory"
LLDP_UPDATE_BIN_NAME = "nic-updater"
BENCHMARK_BIN_NAME = "benchmark"
BENCHMARK_SCHEDULER_BIN_NAME = "benchmark-scheduler"

DOCKER_REGISTRY = "localhost:5000"
DOCKER_IMAGE_NAME = "inventory"
DOCKER_IMAGE_TAG = "latest"

GOPRIVATE ?= "github.com/onmetal/*"

.PHONY: all
all: compile

.PHONY: compile
compile: fmt vet
	for BIN_NAME in $(INVENTORY_BIN_NAME) $(LLDP_UPDATE_BIN_NAME) $(BENCHMARK_BIN_NAME) $(BENCHMARK_SCHEDULER_BIN_NAME); do \
		go build -o dist/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
	done
	cp -rf res/ dist/

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: dl-pciids
dl-pciids:
	curl https://pci-ids.ucw.cz/v2.2/pci.ids --output ./res/pci.ids

.PHONY: docker-build
docker-build:
	docker build . -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --build-arg GOPRIVATE=$(GOPRIVATE) --build-arg GIT_USER=$(GIT_USER) --build-arg GIT_PASSWORD=$(GIT_PASSWORD)

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: clean
clean:
	rm -rf ./dist/

test:
	@echo "--> Project testing"
	go test -race -v ./... -coverprofile cover.out
	@echo "--> Making coverage html page"
	go tool cover -html cover.out -o ./index.html
