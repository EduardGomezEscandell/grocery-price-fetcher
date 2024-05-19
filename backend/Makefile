.PHONY: tidy lint

VERSION ?= "dev"
GO := go build -ldflags -X=github.com/EduardGomezEscandell/grocery-price-fetcher/pkg/services/version.Version=$(VERSION)

tidy: ## Update and clean Go dependencies
	go mod tidy

lint: ## Run linter (installs golangci-lint if not found)
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

build: tidy ## Build Go binaries
	mkdir -p ../bin
	$(GO) -o ../bin/compra cmd/compra/main.go	
	$(GO) -o ../bin/grocery-server cmd/server/main.go
