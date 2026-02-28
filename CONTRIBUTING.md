# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [Git](https://git-scm.com/) 2.0+

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the application:
   ```bash
   go run ./cmd/main.go
   ```
5. Build the project:
   ```bash
   go build -o bin/investmate ./cmd/main.go
   ```
6. Run tests:
   ```bash
   go test ./...
   ```
7. Run tests with coverage:
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func coverage.out
   ```
8. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow)
9. Open a pull request against `main`
