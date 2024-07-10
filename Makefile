build:
	@go build -o bin/api-notez

run: 
	@go run ./cmd/main.go

test: 
	@go test -v ./...

