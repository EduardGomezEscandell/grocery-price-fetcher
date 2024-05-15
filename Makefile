.PHONY: tidy build-go build-js lint test update-golden quality run-mock containerize push run stop clean

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
	cd frontend && npm install && npm run build

run-mock: stop
# Serves the frontend with a mock back-end
# Fast to spin up
	cd frontend && npm run start

containerize: build-go build-js
	cd deploy/container && ./filesystem.sh build
	cd deploy/container && sudo docker build . -t grocery-price-fetcher

push: containerize
	sudo docker tag grocery-price-fetcher edugomez/grocery-price-fetcher:latest
	sudo docker push edugomez/grocery-price-fetcher:latest

run: stop containerize
	cd deploy && sudo docker-compose up

stop:
	sudo docker container rm -f `sudo docker container ls -a | grep grocery-server | cut -c-12` || true

clean: stop
	rm -r bin
	cd deploy && ./filesystem.sh clean