.PHONY: tidy build-go build-js lint test update-golden quality run-mock build-docker push start stop clean

tidy:
	go mod tidy

build-go: tidy
	mkdir -p bin
	go build -o bin/compra cmd/compra/main.go
	go build -o bin/grocery-server cmd/server/main.go

lint:
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

test: build-go
	go test ./...

update-golden: build-go
	UPDATE_GOLDEN=1 go test ./...

quality: build-go lint test

build-js:
	cd frontend && npm install
	cd frontend && npm run build

run-mock: stop
# Serves the frontend with a mock back-end
# Fast to spin up
	cd frontend && npm run start

build-docker: build-go build-js
	cd deploy/container && make build

package: push
	cd deploy/host && make package

full-start: build-docker start

start:
	cd deploy/host && make start

stop:
	cd deploy/host && make stop

clean:
	rm -r bin || true
	cd deploy/container && make clean
	cd deploy/host && make purge