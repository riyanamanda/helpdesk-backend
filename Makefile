.PHONY: build run run-worker tidy test clean seed migrate-up migrate-down migrate-create format

build:
	go build ./...

run:
	@trap 'kill 0' EXIT; \
	air & \
	air -c .air.worker.toml

tidy:
	go mod tidy
	
test:
	@echo "Running tests..."
	go test -v ./...
	
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@go clean

seed:
	go run ./cmd/seed

migrate-up:
	goose up

migrate-down:
	goose down

migrate-create:
	goose create -s $(name) sql

format:
	goimports -w -local github.com/riyanamanda/helpdesk-backend .