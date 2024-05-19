.PHONY: tidy build-go build-js lint test update-golden quality run-mock build-docker package deploy full-start install start stop clean


VERSION := $(shell git describe --tags --always --dirty)
GO := go build -ldflags -X=github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/version.Version=$(VERSION)

help: ## Show this help message
	@echo "Grocery Price Fetcher Development Makefile"
	@echo ""
	@echo "COMMANDS"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tidy: ## Update and clean Go dependencies
	go mod tidy

build-go: tidy ## Build Go binaries
	mkdir -p bin
	$(GO) -o bin/compra cmd/compra/main.go	
	$(GO) -o bin/grocery-server cmd/server/main.go

lint: ## Run linter (installs golangci-lint if not found)
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

test: build-go ## Run tests
	go test ./...

update-golden: build-go ## Update golden test files
	UPDATE_GOLDEN=1 go test ./...

quality: build-go lint test ## Run linter, tests, and quality checks

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

clean: ## Clean up build artifacts, Docker containers, and images
	rm -r bin || true
	cd deploy/container && make clean
	cd deploy/host && make purge