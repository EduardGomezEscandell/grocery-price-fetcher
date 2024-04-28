.PHONY: tidy build-go build-js test update-golden lint quality clean containerize run stop

tidy:
	go mod tidy

build-go: tidy
	mkdir -p bin
	go build -o bin/compra cmd/compra/main.go
	go build -o bin/needs cmd/needs/main.go
	go build -o bin/grocery-server cmd/server/main.go

build-js:
	cd frontend && npm install && npm run build

containerize: build-go build-js
	cd deploy && ./filesystem.sh build
	cd deploy && sudo docker build . -t grocery-server

run: stop containerize
	sudo docker run --name grocery-server -p 80:3000 docker.io/library/grocery-server

stop:
	sudo docker container rm -f `sudo docker container ls -a | grep grocery-server | cut -c-12` || true

test: build-go
	go test ./...

update-golden: build-go
	UPDATE_GOLDEN=1 go test ./...

lint:
	$$(go env GOPATH)/bin/golangci-lint version \
		|| curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	$$(go env GOPATH)/bin/golangci-lint run ./...

quality: build-go lint test

clean: stop
	rm -r bin
	cd deploy && ./filesystem.sh clean