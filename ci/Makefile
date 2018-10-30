.PHONY: default build clean test-all test-unit test-integration test-cover

default: build

# Build binaries
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./cmd/worker/main ./cmd/worker/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./cmd/producer/main ./cmd/producer/*.go
	@echo "Built."

# Remove binaries
clean:
	rm ./cmd/worker/main
	rm ./cmd/producer/main
	@echo "Cleaned."

# Run unit and integration tests
test-all: test-unit test-integration
	@echo "Finished."

# Run unit tests
test-unit:
	go test -tags=unit ./...

# Run integration tests
test-integration:
	docker run -d --rm -p 5672:5672 --name rabbitmq rabbitmq
	go test -tags=integration ./...
	docker stop rabbitmq

# Get unit test coverage
test-cover:
	go test -coverprofile=coverage.out -tags=unit ./... > /dev/null 
	go tool cover -func=coverage.out 
	rm coverage.out
