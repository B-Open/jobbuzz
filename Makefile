build:
	go build ./...

test:
	go test ./...

generate:
	go run github.com/99designs/gqlgen generate
