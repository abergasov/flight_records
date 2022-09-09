APP_PORT := $(or ${port},${port},"8080")

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Check code style
	golangci-lint run --verbose

run: ## Run app. default port is `8080`. You can change it with `make port=8081 run`
	@echo "Running on port ${APP_PORT}"
	go mod vendor
	go run ./cmd/flight_records.go -port=${APP_PORT}

test: ## Run tests
	go test -race ./...

benc: ## Run benchmarks
	go test -bench=. ./...