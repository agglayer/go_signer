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
GOBINARY := aggkit
GOCMD := $(GOBASE)/cmd


# Check dependencies
# Check for Go
.PHONY: check-go
check-go:
	@which go > /dev/null || (echo "Error: Go is not installed" && exit 1)
	
	
# Targets that require the checks
build: check-go
lint: check-go

.PHONY: build-aggkit
build-aggkit:
	$(GOENVVARS) go build -ldflags "all=$(LDFLAGS)" -o $(GOBIN)/$(GOBINARY) $(GOCMD)
	
.PHONY: test-unit
test-unit:
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -short -race -p 1 -covermode=atomic -coverprofile=coverage.out  -coverpkg ./... -timeout 15m ./...
	
.PHONY: lint
lint: ## Runs the linter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/golangci-lint run --timeout 5m


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
