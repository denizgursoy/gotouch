
build:
	go build -o gotouch-mac -v ./cmd/gotouch/

build-test:
	env GOOS=linux GARCH=amd64 go build --ldflags="-X 'github.com/denizgursoy/gotouch/internal/prompts.Environment=test'" -v -o gotouch-linux-test ./cmd/gotouch/

unit-test:
	go test -v ./...

integration-test:
	go test -v --tags=integration ./...