.PHONY: help tidy lint build test coverage

VERSION ?= "dev"
GO := go build -ldflags -X=github.com/EduardGomezEscandell/grocery-price-fetcher/backend/pkg/services/version.Version=$(VERSION)
BIN := ../bin
COVDIR := $(BIN)/coverage

help: ## Show this help message
	@echo "Grocery Price Fetcher Development Makefile"
	@echo ""
	@echo "COMMANDS"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tidy: ## Update and clean Go dependencies
	go mod tidy

lint: ## Run linter (installs golangci-lint if not found)
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

build: ## Build Go binaries
	mkdir -p $(BIN)
	$(GO) -o $(BIN)/grocery-server cmd/server/main.go

test: ## Run unit tests
	mkdir -p $(COVDIR)
	go test ./... -shuffle=on -race

update-golden: export UPDATE_GOLDEN = 1
update-golden: test ## Update golden files

coverage: ## Show test coverage
	mkdir -p $(COVDIR)
	go test ./... -coverprofile=$(COVDIR)/.coverage -covermode=set
	go tool cover -html=$(COVDIR)/.coverage -o $(COVDIR)/index.html \
		&& python3 -m http.server 8080 --directory $(COVDIR)
