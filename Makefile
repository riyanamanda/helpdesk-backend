.PHONY: run run-worker tidy test clean seed migrate-up migrate-down format

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

format:
	goimports -w -local github.com/riyanamanda/helpdesk-backend .