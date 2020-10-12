SHELL := /bin/bash
GO111MODULE := on

.PHONY: all
all: help

.PHONY: test
test: ## Run unittests
	@go test  ./... -v

.PHONY: build
build: ## build the current go project
	cd ./cmd/gbm && go build -v main.go && cd ../..

.PHONY: run
run: ## run the current go project
	cd ./cmd/gbm && go run -v main.go && cd ../..


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

