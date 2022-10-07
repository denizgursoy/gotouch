LINUX_BUILD := env GOOS=linux GARCH=amd64 go build

build:
	go build -o gotouch-mac -v ./cmd/gotouch/

fmt:
	go fmt ./...

build-test:
	 $(LINUX_BUILD) -v -o gotouch-linux-test -tags=integration ./cmd/gotouch/

build-mac-test:
	go build -v -o gotouch-mac-test ./cmd/gotouch/

unit-test:
	go test -v --tags=unit ./...

integration-test: build-test
	go test -v --tags=integration_test ./...

generate-mocks:
	go generate -v ./...