APP=sftd
PKG=./...
GOFLAGS?=
LDFLAGS?=-s -w

.PHONY: run test lint fmt

run:
	go run $(GOFLAGS) ./cmd/sftd

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w . && go vet $(PKG)
