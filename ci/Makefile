
tags=""


all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./worker/main ./worker/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o ./producer/main ./producer/*.go
	@echo "Built."

clean:
	rm ./worker/main
	rm ./producer/main
	@echo "Cleaned."

test:
	docker run -d --rm -p 5672:5672 --name rabbitmq rabbitmq
	go test ./...
	docker stop rabbitmq

cover:
	docker run -d --rm -p 5672:5672 --name rabbitmq rabbitmq
	go test -coverprofile=coverage.out -tags=$(tags) ./... > /dev/null 
	go tool cover -func=coverage.out 
	rm coverage.out
	docker stop rabbitmq

