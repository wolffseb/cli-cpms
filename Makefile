BINARY := cpms
PKG     := github.com/wolffseb/cli-cpms

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(PKG)/internal/buildinfo.Version=$(VERSION) \
	-X $(PKG)/internal/buildinfo.Commit=$(COMMIT) \
	-X $(PKG)/internal/buildinfo.Date=$(DATE)

.PHONY: all build test race cover vet lint fmt check-fmt tidy clean ci

all: build

## build: compile the cpms binary into bin/
build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/cpms

## test: run the test suite
test:
	go test ./...

## race: run the test suite under the race detector
race:
	go test -race ./...

## cover: run tests and open a coverage summary
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

vet:
	go vet ./...

## lint: run golangci-lint (go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2)
lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

check-fmt:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then echo "not gofmt'd:"; echo "$$out"; exit 1; fi

tidy:
	go mod tidy

clean:
	rm -rf bin coverage.out

## ci: everything the CI workflow runs, locally
ci: check-fmt vet lint race
