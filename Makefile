LINUX_BUILD := env GOOS=linux go build
MAC_BUILD := env GOOS=darwin go build

.PHONY: build build-test fmt lint unit coverage integration mock release install-semver up-major up-minor up-release

build:
	go build -o gotouch -v ./cmd/gotouch/

fmt:
	go fmt ./...
	gofumpt -l -w  -extra .

lint:
	golangci-lint run

build-test:
	$(LINUX_BUILD) -v -o gotouch-linux -tags=integration ./cmd/gotouch/
	$(MAC_BUILD)   -v -o gotouch-darwin -tags=integration ./cmd/gotouch/

unit:
	go test -v -coverprofile coverage.out ./...

coverage:
	go tool cover -html=coverage.out

integration: build-test
	go test -v --tags=integration_test ./...

mock:
	go generate -v ./...

release:
	goreleaser release --snapshot --clean

install-semver:
	go install github.com/maykonlf/semver-cli/cmd/semver@latest

up-major:
	semver up major

up-minor:
	semver up minor

up-release:
	semver up release