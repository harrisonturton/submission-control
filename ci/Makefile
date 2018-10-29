
all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./worker/main ./worker/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./producer/main ./producer/*.go
	@echo "Built."

clean:
	rm ./worker/main
	rm ./producer/main
	@echo "Cleaned."

test:
	rabbitmq-server &
	go test ./...

cover:
	go test -coverprofile=coverage.out ./... > /dev/null 
	go tool cover -func=coverage.out 
	rm coverage.out

