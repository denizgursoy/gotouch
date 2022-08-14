
all:
	go build -v ./cmd/gotouch/

unit-test:
	go test -v  --tags=unit ./...


integration-test:
	go test -v  --tags=integration ./...