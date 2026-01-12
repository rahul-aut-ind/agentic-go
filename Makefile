run-local:
	go run cmd/service-genkit/main.go	

sanitize: lint test

lint:
	golangci-lint run --timeout=2m

test:
	go test ./... -v
	