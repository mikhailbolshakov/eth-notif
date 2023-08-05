.PHONY: dep lint build vendor

# load env variables from .env
ENV_PATH ?= ./.env
ifneq ($(wildcard $(ENV_PATH)),)
    include .env
    export
endif

export GOFLAGS=-mod=vendor

# Build commands =======================================================================================================

vendor:
	go mod vendor

dep:
	go mod tidy

check-lint-installed:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed"; \
		exit 1; \
	fi; \

lint: check-lint-installed
	@echo Running golangci-lint
	golangci-lint --modules-download-mode vendor run --timeout 5m0s --skip-dirs-use-default ./...
	go fmt ./...

build: lint vendor ## builds the main
	@mkdir -p bin
	go build -o bin/ cmd/main.go

run: ## run the service
	./bin/main

# Tests commands =======================================================================================================

test: ## run the tests
	go test -count=1 ./...


