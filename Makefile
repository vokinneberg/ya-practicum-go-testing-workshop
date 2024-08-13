.PHONY: build test test-coveralls deploy diff destroy gen lint config setup stop teardown

# Export envars
include .env
export

define build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o bin/$(1) cmd/$(1)/main.go
endef

define run
	go run ./cmd/$(1)/main.go
endef

build:
	$(call build,shortener)

run:
	$(call run,shortener)

test:
	go test ./...

test-coveralls:
	go test -v -cover -coverprofile=profile.cov ./...

# code generation for mock
gen:
	GOFLAGS=-mod=mod go generate ./...

lint:
	golangci-lint run

# starts dependent services inside docker containers
setup:
	docker compose up -d --remove-orphans --build
	docker compose run wait -c postgres:5432 -t 80

# stops containers
stop:
	docker compose stop

# stops and removes all dependent service containers
teardown:
	docker compose down --remove-orphans

