
all:
	go build -o gotouch-mac -v ./cmd/gotouch/
	env GOOS=linux GARCH=amd64 go build -v -o gotouch-linux ./cmd/gotouch/
	go build --ldflags="-X 'github.com/denizgursoy/gotouch/internal/prompts.Environment=test'" -o gotouch-mac-test -v ./cmd/gotouch/
	env GOOS=linux GARCH=amd64 go build --ldflags="-X 'github.com/denizgursoy/gotouch/internal/prompts.Environment=test'" -v -o gotouch-linux-test ./cmd/gotouch/

unit-test:
	go test -v --tags=unit ./...


integration-test:
	go test -v --tags=integration ./...