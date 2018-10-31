.PHONY: default build start clean test test-unit test-integration test-cover

# Configuration to run the multi-container service
CONFIG = docker-compose.yml
TEST_CONFIG = docker-compose.test.yml

default: build

# Build binaries
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./cmd/worker/main ./cmd/worker/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./cmd/producer/main ./cmd/producer/*.go
	@echo "Built."

# Start the service
start:
	docker-compose -f $(CONFIG) up

# Remove binaries
clean:
	rm ./cmd/worker/main
	rm ./cmd/producer/main
	@echo "Cleaned."

# Run unit and integration tests
test:
	@echo "\033[4mLaunching test containers...\033[0m"
	docker-compose --project-name test-integration -f $(TEST_CONFIG) up -d
	-go test -tags="unit integration" ./...
	@echo "\033[4mStopping test containers. This may take a few seconds...\033[0m"
	docker-compose --project-name test-integration down

# Run unit tests
test-unit:
	go test -tags=unit ./...

# Run integration tests
test-integration:
	@echo "\033[4mLaunching test containers...\033[0m"
	docker-compose --project-name test-integration -f $(TEST_CONFIG) up -d
	-go test -tags=integration ./...
	@echo "\033[4mStopping test containers. This may take a few seconds...\033[0m"
	docker-compose --project-name test-integration down

# Unit test coverage
test-cover:
	go test -coverprofile=coverage.out -tags=unit ./... > /dev/null 
	go tool cover -func=coverage.out 
	rm coverage.out
