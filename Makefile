APP_NAME := itibar-scraper
VERSION := 1.0.0

default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

vet: ## runs go vet
	go vet ./...

format: ## runs go fmt
	gofmt -s -w .

test: ## runs the unit tests
	go test -v -race -timeout 5m ./...

test-cover: ## outputs the coverage statistics
	go test -v -race -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	rm coverage.out

lint: ## runs the linter
	go tool golangci-lint -v run ./...

build: ## builds the application (default: playwright)
	go build -o bin/$(APP_NAME) .

build-rod: ## builds the application with go-rod browser engine
	go build -tags rod -o bin/$(APP_NAME)-rod .

docker: ## builds docker image with playwright (default)
	docker build -t $(APP_NAME):$(VERSION) .

docker-rod: ## builds docker image with go-rod
	docker build -f Dockerfile.rod -t $(APP_NAME):$(VERSION)-rod .

clean: ## clean build artifacts
	@rm -rf bin/ tmp/
