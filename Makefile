LINUX_BUILD := env GOOS=linux GARCH=amd64 go build

build:
	go build -o gotouch-mac -v ./cmd/gotouch/

build-test:
	 $(LINUX_BUILD) --ldflags="-X 'github.com/denizgursoy/gotouch/internal/manager.Environment=test'" -v -o gotouch-linux-test ./cmd/gotouch/

build-mac-test:
	go build --ldflags="-X 'github.com/denizgursoy/gotouch/internal/manager.Environment=test'" -v -o gotouch-mac-test ./cmd/gotouch/

unit-test:
	go test -v --tags=unit ./...

integration-test: build-test
	go test -v --tags=integration ./...

generate-mocks:
	go generate -v ./...