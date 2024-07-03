.PHONY: help build-go test-go test-e2e update-golden update-golden build-js run-mock build-docker package deploy full-start install start stop uninstall clean

VERSION := $(shell git describe --tags --always --dirty)

help: ## Show this help message
	@echo "Grocery Price Fetcher Development Makefile"
	@echo ""
	@echo "COMMANDS"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-go: ## Build Go binaries
	cd backend && make tidy
	cd backend && make lint
	cd backend && make build VERSION=$(VERSION)

test-go: build-go ## Run unit tests
	cd backend && make test

test-e2e: build-go ## Run end-to-end tests
	cd end-to-end && go test ./... -count=1 -race -shuffle on

update-golden: export UPDATE_GOLDEN = 1
update-golden: test-e2e ## Update golden test files

build-js: ## Build the frontend JavaScript code
	cd frontend && npm install
	cd frontend && npm run build

run-mock: ## Start the frontend with a mock backend
# Serves the frontend with a mock back-end
# Fast to spin up
	cd frontend && npm run start

build-docker: build-go build-js ## Build the Docker image
	cd deploy/container && make build

package: ## Package the application for deployment (see deploy/host Makefile)
	cd deploy/host && make package

deploy: ## Deploy the application (see deploy/host Makefile)
	cd deploy/host && make deploy

full-start: build-docker install start ## Build the application and self-host it

FQDN ?= https://localhost
install: ## Install the the application locally (see deploy/install Makefile)
	cd deploy/host && make install FQDN=$(FQDN)

start: ## Self-host the application (see deploy/start Makefile)
	cd deploy/host && make start

stop: ## Stop the application (see deploy/host Makefile)
	cd deploy/host && make stop

uninstall: ## Uninstall the application (see deploy/purge Makefile)
	cd deploy/host && make purge

clean: ## Clean up build artifacts, Docker containers, and images
	rm -r bin || true
	cd deploy/container && make clean
	cd deploy/host && make purge
	cd backend && make clean