PROVIDER := pulumi-resource-hyperv
VERSION := $(shell pulumictl get version)
TESTPARALLELISM := 4
WORKING_DIR := $(shell pwd)
SUB_PROJECTS := provider examples
GOOS := windows

.PHONY: provider build_sdks build clean test lint provider_debug

build:: provider

provider::
	cd provider && GOOS=$(GOOS) go build -o bin/$(PROVIDER) ./cmd/$(PROVIDER)

provider_debug::
	cd provider && GOOS=$(GOOS) go build -gcflags="all=-N -l" -o bin/$(PROVIDER) ./cmd/$(PROVIDER)

lint::
	cd provider && GOOS=$(GOOS) golangci-lint run --path-prefix=provider -c ../.golangci.yml
	npx markdownlint "**/*.md"

test::
	cd provider && GOOS=$(GOOS) go test -v ./...

test_provider::
	cd provider && GOOS=$(GOOS) go test -v ./...

test_examples::
	cd examples && go test -v ./...

format::
	gofmt -w .

clean::
	rm -rf $(WORKING_DIR)/bin
	cd provider && rm -rf bin
	cd sdk && rm -rf bin
	cd $(WORKING_DIR)/sdk/nodejs && rm -rf node_modules/

help::
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build            Build the provider and SDKs"
	@echo "  provider         Build the provider only"
	@echo "  provider_debug   Build the provider with debug symbols"
	@echo "  build_sdks       Build the SDKs"
	@echo "  test             Run tests"
	@echo "  test_provider    Run provider tests"
	@echo "  test_examples    Run example tests"
	@echo "  lint             Run linters"
	@echo "  format           Format code"
	@echo "  clean            Clean build artifacts"