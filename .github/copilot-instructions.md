# Copilot Instructions for InvestMate

## Project Overview

InvestMate is a Go CLI application that fetches, processes, and displays ETF (Exchange-Traded Fund)
financial data. It retrieves dividend cash amounts and average closing prices from external sources,
calculates dividend yields per year, and renders the results as a colour-coded ASCII table in the
terminal.

## Repository Structure

```
investmate/
├── cmd/
│   ├── main.go              # Entry point: orchestrates ETF processing and table rendering
│   └── main_test.go         # Tests for main-package functions
├── internal/
│   ├── domain/
│   │   └── entities/
│   │       ├── etf.go       # ETF struct and all calculation/formatting methods
│   │       └── etf_test.go  # Unit tests for entity logic
│   └── infrastructure/
│       └── repositories/
│           ├── nasdaq/      # NASDAQ REST API adapters (dividends + closing prices)
│           ├── statusinvest/# StatusInvest web-crawler adapter (dividends, currently unused)
│           └── historyorg/  # History.org crawler adapter (currently unused)
├── .github/
│   └── workflows/
│       └── default.yaml     # CI/CD pipeline (delegates to shared reusable workflow)
├── go.mod / go.sum          # Go module definition and checksums
├── horusec.json             # Horusec SAST configuration
├── CHANGELOG.md             # Keep a Changelog format with semantic versioning
├── CONTRIBUTING.md          # Development prerequisites and workflow
└── README.md                # Project description and feature overview
```

## Technology Stack & Dependencies

| Dependency | Purpose |
|---|---|
| `github.com/gocolly/colly` | HTML scraping framework |
| `github.com/olekukonko/tablewriter` | Renders ASCII tables with colour support |
| `github.com/sirupsen/logrus` | Structured, levelled logging |
| `github.com/stretchr/testify` | Test assertions (`assert` package) |

Go version: **1.26+** (declared in `go.mod`).

## Architecture & Design Patterns

- **Clean Architecture** with a domain layer (`internal/domain`) and infrastructure layer
  (`internal/infrastructure`). Domain entities have no external dependencies.
- **Repository pattern** — each data source is encapsulated behind its own repository struct. New
  sources can be added under `internal/infrastructure/repositories/` without touching domain code.
- The `cmd/main.go` orchestration layer wires repositories to domain entities and drives rendering.

## Build, Test, Lint & Run Commands

```bash
# Install / tidy dependencies
go mod download
go mod tidy

# Run the application (fetches live data from NASDAQ)
go run ./cmd/main.go

# Build a binary
go build -o bin/investmate ./cmd/main.go

# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func coverage.out

# Lint (golangci-lint must be installed)
golangci-lint run ./...
```

> **Timings:** `go test ./...` completes in under 30 seconds on a standard machine when network
> access is available (some tests call live APIs). `go build` completes in under 10 seconds.

## CI/CD Pipeline

CI is defined in `.github/workflows/default.yaml` and delegates entirely to the shared reusable
workflow at `rios0rios0/pipelines/.github/workflows/go-binary.yaml@main`. The pipeline runs on:

- every push to `main`
- every tag push
- every pull request targeting `main`
- manual `workflow_dispatch`

The pipeline produces a binary named `investmate` and publishes releases automatically on tagged
commits. SonarCloud is used for static analysis and coverage reporting.

## Coding Conventions

- Package names are lowercase, single words (e.g. `nasdaq`, `entities`).
- Exported identifiers use PascalCase; unexported identifiers use camelCase.
- Error wrapping: use `fmt.Errorf("context: %w", err)`.
- Logging: use `logrus` with structured fields (`logger.WithError(err).Errorf(...)`).
- Tests follow **Arrange / Act / Assert** (`// given`, `// when`, `// then` comments) and use
  `t.Parallel()` at both the parent and sub-test levels.
- Test functions are named `Test<Package>_<Function>` with descriptive sub-test names in the form
  `"should … when …"`.
- Constants are UPPER_SNAKE_CASE (e.g. `YearsToFetch`, `PercentageMultiplier`).
- Commits follow [Conventional Commits](https://www.conventionalcommits.org/) and the project's
  [Git Flow guide](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow).

## Configuration

The two key compile-time constants in `cmd/main.go` control behaviour:

```go
const (
    YearsToFetch = 5 // How many years of historical data to fetch and display
)
```

The list of ETFs to process is defined as a slice in `main()`:

```go
etfNames := []string{
    "SPY", "QQQ", "SCHD", "YYY", "GLD",
    "HYGW", "RIET", "SDIV", "SVOL", "XYLD",
}
```

## Development Workflow

1. Fork the repository and create a feature branch: `git checkout -b feat/my-change`
2. Install dependencies: `go mod download`
3. Make changes; add or update tests under the same package as the code being changed.
4. Run `go test ./...` to verify all tests pass.
5. Run `golangci-lint run ./...` to check for lint issues.
6. Commit using Conventional Commits and open a pull request against `main`.

## Common Tasks

| Task | Command |
|---|---|
| Add a new ETF source | Create a new package under `internal/infrastructure/repositories/`, implement `ListDividendsByETF` and/or `ListClosingPricesByETF`, then wire it in `cmd/main.go` |
| Add a new ETF ticker | Append its symbol to `etfNames` in `cmd/main.go` |
| Increase historical range | Change `YearsToFetch` constant in `cmd/main.go` |
| Check SAST findings | Review `horusec.json` and run `horusec start` |

## Troubleshooting

- **Empty data for an ETF** — the NASDAQ API blocks requests without a browser-like `User-Agent`.
  The repository sets one explicitly; verify it has not changed upstream.
- **Build errors after `go mod tidy`** — ensure your local Go version is ≥ 1.26.
- **Table colours not showing** — some terminals do not support ANSI colour codes; run in a
  terminal that does (e.g. `bash`, `zsh`, Windows Terminal).
