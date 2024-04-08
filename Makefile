tidy:
	go mod tidy

build: tidy
	mkdir -p bin
	go build -o bin/compra cmd/compra/main.go

test: build
	go test ./...

update-golden: build
	UPDATE_GOLDEN=1 go test ./...

lint:
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

quality: build lint test

clean:
	rm -r bin