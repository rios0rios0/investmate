# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

InvestMate is a Go CLI application that fetches ETF financial data (dividend cash amounts, average closing prices) from external sources, calculates dividend yields per year, and renders results as a colour-coded ASCII table in the terminal.

## Build, Test, and Lint

```bash
# Build / run — entry point is the package at ./cmd (cmd/main.go).
# NOTE: the Makefile build/run/debug/install targets point at ./cmd/investmate,
# which no longer exists, so invoke go directly:
go build -o bin/investmate ./cmd   # Binary to bin/investmate
go run ./cmd                        # Run directly (fetches live NASDAQ data)

# Test
go test ./...             # All tests (some call live NASDAQ APIs, requires network)
go test -race ./...       # With race detector
go test -run TestETFTestSuite ./internal/domain/entities/...  # Single suite
go test -run TestMain_ProcessETF ./cmd/...                    # Single test

# Quality gates (always use Makefile targets, never call tools directly)
make lint                 # Linting via pipelines repo
make sast                 # Full SAST suite (CodeQL, Semgrep, Trivy, Hadolint, Gitleaks)
```

## Architecture

Clean Architecture with Hexagonal (Ports & Adapters) pattern:

- **`cmd/main.go`** — Entry point. Wires repositories to entities, drives table rendering. Manual DI (no framework).
- **`internal/domain/`** — Contracts layer. Contains entity business logic and repository interfaces (ports).
- **`internal/infrastructure/`** — Implementations layer. Contains repository adapters for external data sources.

### Data Flow

`main()` → instantiates repository adapters → calls repository interfaces → populates ETF entity → entity calculates yields → renders ASCII table.

### Repository Interfaces (Ports)

- `DividendsRepository.ListDividendsByETF(etf string) (map[string]float64, error)`
- `PricesRepository.ListClosingPricesByETF(etf string) (map[string]float64, error)`

### Adapters

- **`nasdaq/`** — Active. REST API adapters for dividends and closing prices.
- **`statusinvest/`** and **`historyorg/`** — Currently unused web-crawler adapters.

### Key Dependencies

| Package | Purpose |
|---|---|
| `gocolly/colly` | HTML scraping |
| `olekukonko/tablewriter` | ASCII table rendering with ANSI colours |
| `sirupsen/logrus` | Structured logging |
| `stretchr/testify` | Test assertions and suites |

## Testing Conventions

- Tests use **testify** (`assert` for unit tests, `suite.Suite` for entity tests).
- BDD structure with `// given`, `// when`, `// then` comments.
- Sub-test names follow `"should ... when ..."` pattern.
- Unit tests use `t.Parallel()` at both parent and sub-test levels.
- Entity tests use `suite.Run()` (not parallel due to shared setup).
- Some tests call live NASDAQ APIs — expect network access requirement.

## Configuration

Compile-time constants in `cmd/main.go`:
- `YearsToFetch` (default: 5) — years of historical data
- `targetYieldPercentage` (default: 9) — minimum yield for green colouring
- `etfNames` slice in `main()` — list of ETF tickers to process

## Troubleshooting

- **Empty API data**: NASDAQ blocks requests without a browser-like `User-Agent` header. The repository sets one explicitly.
- **Build errors after `go mod tidy`**: Requires Go >= 1.27.

## CI/CD

Defined in `.github/workflows/default.yaml`, delegates to `rios0rios0/pipelines/.github/workflows/go-binary.yaml@main`. Triggers on push to main, tag push, PR to main, and manual dispatch.

<!-- chlog:start -->
## Changelog (chlog) — MANDATORY

If the repository you are working in uses chlog (a `.chlog.yaml` or `.chlog.yml`
config file, or a `.changes/` directory, exists at the project root), the
following is binding and ALWAYS applies: whenever you make ANY change, you MUST
create a changelog fragment as part of the same change — automatically, without
being asked, before committing.

- Do NOT edit CHANGELOG.md directly; it is generated from fragments.
- Create the fragment with:
  `chlog new --kind <Kind> --body "<imperative description>"`
- Valid kinds: Added, Changed, Deprecated, Removed, Fixed, Security
- Choose the kind that best matches the change (e.g., new feature → Added,
  bug fix → Fixed, behavior change → Changed, removal → Removed, security fix → Security).
- If the change is backward-INCOMPATIBLE with the public API (a breaking
  change), you MUST add the `--breaking` flag:
  `chlog new --kind <Kind> --breaking --body "<description>"`.
  This is the ONLY thing that triggers a major version bump — the kind alone
  never does (per SemVer, major = incompatible change). When unsure whether a
  change breaks compatibility, ask the user instead of guessing.
- Fragments are YAML files in `.changes/unreleased/`; stage them with your commit.
- `chlog check` fails the build when a fragment is missing — never skip it.
<!-- chlog:end -->
