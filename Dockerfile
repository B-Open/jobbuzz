FROM golang:1.17 AS build

# Download packages
WORKDIR /go/src/github.com/b-open/jobbuzz
COPY go.* ./
RUN go mod download

# Build binaries
COPY ./ ./
RUN mkdir -p bin
RUN go build -o ./bin ./...

# Make final image
FROM debian:11

WORKDIR /app
COPY --from=build /go/src/github.com/b-open/jobbuzz/bin/db-migrator ./
COPY --from=build /go/src/github.com/b-open/jobbuzz/bin/jobbuzz-api ./
COPY --from=build /go/src/github.com/b-open/jobbuzz/bin/jobbuzz-scraper ./

ENTRYPOINT ["/app/db-migrator"]
