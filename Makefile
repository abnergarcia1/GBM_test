SHELL := /bin/bash
##IMPORT_CONTAINER=voxie/engineering-test:incoming
GO111MODULE := on

.PHONY: all
all: help

.PHONY: redis
import: ## Run redis cointainer
	@echo "+ $@"
	@docker run -p 6379:6379 --name jaca-core-redis -d redis

.PHONY: build
build: ## build the current go project
	cd ./cmd/gbm && go build -v main.go && cd ../..

.PHONY: run
run: ## run the current go project
	cd ./cmd/gbm && go run -v main.go && cd ../..

.PHONY: create_docker_local
create_docker_local:  ## Build docker images for deploy
	docker build -t jaca-restapi -f Dockerfile .


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

