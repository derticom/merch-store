lint:
	golangci-lint run ./... -v

test:
	go test ./internal/...