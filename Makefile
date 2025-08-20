build:
	@go build -o bin/catalog_service

run: build
	@./bin/catalog_service

test:
	@go test -v ./...
