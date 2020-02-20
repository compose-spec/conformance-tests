.DEFAULT_GOAL := help

.PHONY: images
images: ## Build the test images
	docker build server -t test-server

.PHONY: test
test: images ## Run tests
	GOPRIVATE=github.com/compose-spec/compatibility-test-suite go test ./... -v

.PHONY: fmt
fmt: ## Format go files
	@goimports -e -w ./

.PHONY: lint
lint: ## Verify Go files
	golangci-lint run --config ./golangci.yml ./

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
