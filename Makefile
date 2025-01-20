lint:
	@golangci-lint run ./...
	@echo "ok\t\tlinter passed"

test:
	@go test ./... -race -coverprofile=coverage.out
	@echo "ok\t\tunit tests passed"

cover: test
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ "$$coverage" != "100.0" ]; then \
		echo "fail\ttest coverage is insufficient"; \
		exit 1; \
	else \
		echo "ok\t\ttest coverage passed"; \
	fi

pre-commit: test cover lint
	@echo "\nAll pre-commit checks have passed!"