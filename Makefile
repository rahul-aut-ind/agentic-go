run-local:
	go run cmd/service-genkit/main.go	

sanilize: deps lint

lint:
	golangci-lint run --timeout=2m

deps:
	wire ./...