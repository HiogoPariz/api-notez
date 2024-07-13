build:
	@go build -o bin/api-notez

run: 
	@go run ./cmd/main.go

test: 
	@go test -v ./...

migration:
	@read -p "Migration name: " migration; \
	GOOSE_DRIVER=postgres; \
	GOOSE_DBSTRING="user=postgres dbname=postgres password=api-notez sslmode=disable"; \
	goose -dir="./internal/migration/" create $$migration sql
