build:
	go build -tags rod -o bin/itibar-scraper

test:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w .
