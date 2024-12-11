lint:
	golangci-lint run ./...

test:
	go test ./... -race

pre-commit: lint test
