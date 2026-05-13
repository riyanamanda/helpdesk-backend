.PHONY: run tidy test clean seed migrate-up migrate-down

run:
	air

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