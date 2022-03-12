# JobBuzz

[![Go build and test](https://github.com/B-Open/jobbuzz/actions/workflows/go.yml/badge.svg)](https://github.com/B-Open/jobbuzz/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/B-Open/jobbuzz/branch/main/graph/badge.svg?token=MS1L2JJCD5)](https://codecov.io/gh/B-Open/jobbuzz)
[![Go Report Card](https://goreportcard.com/badge/github.com/b-open/jobbuzz)](https://goreportcard.com/report/github.com/b-open/jobbuzz)

Brunei open source job search database and alert notification

## Development

### Requirements

1. Go 1.17 or higher
2. MySQL 8

### Running locally

#### Natively

1. Copy the `.env.example` file in the repository root to `.env`.
2. Update the contents of `.env` with your database access details.
3. Change working directory to `cmd/jobbuzz-api` and run `go run .` to start the API server.
4. Change working directory to `cmd/jobbuzz-scraper` and run `go run .` to start the web scraper program.

#### Docker

1. Copy the `.env.example` file in the repository root to `.env`.
2. Update the contents of `.env` with your database access details.
3. Run `docker-compose up`.
4. You may need to wait for a while for the scraper to run and complete.
5. The API server can be accessed at http://localhost:8080

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
