.PHONY: build-go build-js test update-golden run-mock build-docker package deploy full-start install start stop uninstall clean

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

test: build-go ## Run tests
	cd end-to-end && go test ./...

update-golden: build-go ## Update golden test files
	UPDATE_GOLDEN=1 go test ./...

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

install: ## Install the the application locally (see deploy/install Makefile)
	cd deploy/host && make install

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