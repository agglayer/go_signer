ARCH := $(shell arch)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
else
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/target
GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=1 GOARCH=$(ARCH)


# Check dependencies
# Check for Go
.PHONY: check-go
check-go:
	@which go > /dev/null || (echo "Error: Go is not installed" && exit 1)
	
	
# Targets that require the checks
build: check-go
lint: check-go

.PHONY: clean
clean:
	rm -rf $(GOBIN)/*


.PHONY: build
build:
	$(GOENVVARS) go build -ldflags "all=$(LDFLAGS)" -o $(GOBIN)/cmdline_signer test/cmdline_signer/main.go
	
.PHONY: test-unit
test-unit:
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -short -race -covermode=atomic -coverprofile=coverage_short.out  -coverpkg ./... -timeout 15m ./...

.PHONY: test-e2e
test-e2e:
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -p 1 -covermode=atomic -coverprofile=coverage.out  -coverpkg ./... -timeout 15m ./...
	

.PHONY: lint
lint: ## Runs the linter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/golangci-lint run --timeout 5m

.PHONY: check-is-new-version
check-is-new-version: ## Checks if the version is new or already exists
	@export VERSION=$$(go run ./cmd/ --version  | cut -f 3 -d ' ') ; \
	echo "current version: $$VERSION" ; \
	if [ -z "$$VERSION" ]; then echo "Error: Version is empty" && exit 1; fi ; \
	git tag -l $$VERSION | grep $$VERSION  ; \
	if [ $$? -eq 0 ]; then echo "Error: Version already exists"; exit 1; fi ; \
	echo "Version is new"

COMMON_MOCKERY_PARAMS=--disable-version-string --with-expecter --exported

.PHONY: generate-mocks
generate-mocks:	
	mockery ${COMMON_MOCKERY_PARAMS}

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
