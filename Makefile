APP=sftd
PKG=./...
GOFLAGS?=
LDFLAGS?=-s -w

.PHONY: dev run build test lint fmt cover

dev: ## quick dev
	@air || reflex -r '\.go$$' -s -- sh -c 'make run' || make run

run:
	go run $(GOFLAGS) ./cmd/sftd

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/sftd

test:
	go test -race -cover -coverprofile=coverage.out $(PKG)

lint:
	golangci-lint run

fmt:
	gofmt -s -w . && go vet $(PKG)

cover:
	go tool cover -func=coverage.out

pprof:
	go tool pprof http://localhost:8080/debug/pprof/profile?seconds=15
