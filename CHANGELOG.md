# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

This file is not edited by hand. Every change writes its own fragment under
`.changes/unreleased/` with [chlog](https://github.com/luizjhonata/chlog), and a release compiles
the pending fragments into a version section here — so two branches each adding an entry no
longer touch the same lines, and a rebase that used to conflict on this file now conflicts on
nothing.

When a new release is proposed:

1. Create a new branch `bump/x.x.x` (this isn't a long-lived branch!!!);
2. The fragments pending under `.changes/unreleased/` are compiled into a version section by `chlog batch auto && chlog merge` (AutoBump does this for you — it reads the fragments directly);
3. Open a Pull Request with the bump version changes targeting the `main` branch;
4. When the Pull Request is merged, a new git tag must be created using [GitHub environment](https://github.com/rios0rios0/investmate/tags).

Releases to productive environments should run from a tagged version.
Exceptions are acceptable depending on the circumstances (critical bug fixes that can be cherry-picked, etc.).

## [Unreleased]

## [0.3.1] - 2026-09-02

### Changed

- changed the Go module dependencies to their latest versions

## [0.3.0] - 2026-08-28

### Added

- added the Claude automated code review and `@claude` mention responder workflows, `claude-review.yaml` and `claude-mention.yaml`, matching the `reusable-claude-review.yaml` / `reusable-claude-mention.yaml` definitions they call in `rios0rios0/pipelines`, authenticating with the `CLAUDE_CODE_OAUTH_TOKEN` secret

### Changed

- changed the Go module dependencies to their latest versions

### Fixed

- restored the `.changes/unreleased/` directory with a `.gitkeep`, so the release tooling keeps recognising this project as [`chlog`](https://github.com/luizjhonata/chlog)-based after a release consumes the last fragment. Git tracks files rather than directories, so the bump commit that removed the final fragment removed the directory too, and the next run read the empty `[Unreleased]` section as "nothing to release"
- restored the `id-token: write` permission on both Claude workflow callers. Without it the caller grants less than the reusable workflow declares, which GitHub rejects before the job starts -- runs ended in `startup_failure`. The action needs the scope because `setupGitHubToken()` exchanges a GitHub OIDC token for the GitHub App token it posts with, unless a `github_token` is passed explicitly.

### Removed

- removed the unused `id-token: write` permission from the Claude workflow callers, and changed `claude-review.yaml`'s display name to `Claude Review` so it matches its file name and its `Claude Mention` sibling. `anthropics/claude-code-action` needs `id-token: write` only for workload identity federation or the Bedrock / Vertex / Foundry OIDC paths; these authenticate with `claude_code_oauth_token`, so the scope allowed minting OIDC tokens for any audience without ever being used.

## [0.2.0] - 2026-08-26

### Added

- added a tailored `code-review` skill under `.github/skills/` so GitHub Copilot reviews changes against the [rios0rios0/guide](https://github.com/rios0rios0/guide/wiki) standards and this repository's own load-bearing invariants

### Changed

- changed the changelog to [chlog](https://github.com/luizjhonata/chlog) fragments: a change now writes its own YAML file under `.changes/unreleased/` through `chlog new --kind <Kind> --body "..."`, and `CHANGELOG.md` is GENERATED from them at release time by `chlog batch auto && chlog merge`. That is the one thing a single shared file cannot do — two branches each adding an entry no longer touch the same lines, so a rebase that used to conflict on `CHANGELOG.md` now conflicts on nothing. The `[Unreleased]` section was empty, so nothing had to be carried across. AutoBump already reads the fragments directly, so the release flow is unchanged.
- changed the Go module dependencies to their latest versions

### Fixed

- fixed the `main` pipeline, which every repository's `sast:gitleaks` job had been failing since the code-review skill landed: the skill's own security bullet listed credential prefixes verbatim to warn against writing them, and the scanner's second pass matches those prefixes on their own, so the warning tripped the rule it was describing. The bullet now names the vendors instead, and the commit that carried the original wording is allowlisted by fingerprint in `.gitleaksignore`, because the scan walks the whole history reachable from `HEAD` and no edit at the tip can clear a past commit. No credential was ever committed.

## [0.1.25] - 2026-08-24

### Changed

- changed struct literals and `errors.As` calls to the Go 1.27 forms required by the `modernize` linter
- changed the Go module dependencies to their latest versions
- changed the Go version to `1.27.0` and updated all module dependencies
- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to reference the Go `1.27` minimum

### Fixed

- fixed `make test` and `make sast` leaving generated reports (`reports/`, `coverage.txt`, `coverage.xml`, `cobertura.xml`, `junit.xml`) as untracked files by adding them to `.gitignore`

## [0.1.24] - 2026-08-17

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.23] - 2026-08-15

### Changed

- changed the Go version to `1.26.6` and updated all module dependencies

## [0.1.22] - 2026-08-13

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.21] - 2026-08-12

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.20] - 2026-08-11

### Changed

- changed the Go module dependencies to their latest versions
- refreshed `.github/copilot-instructions.md` to document the `targetYieldPercentage` compile-time constant

## [0.1.19] - 2026-07-27

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.18] - 2026-07-23

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.17] - 2026-07-22

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.16] - 2026-07-16

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.15] - 2026-07-10

### Changed

- changed the Go module dependencies to their latest versions
- changed the Go version to `1.26.5` and updated all module dependencies

## [0.1.14] - 2026-06-18

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.13] - 2026-06-09

### Changed

- changed the Go module dependencies to their latest versions
- refreshed `CLAUDE.md` build/run commands to invoke `go build`/`go run` against `./cmd` (the Makefile targets still reference the removed `./cmd/investmate` path)

## [0.1.12] - 2026-06-03

### Changed

- changed the Go module dependencies to their latest versions
- changed the Go version to `1.26.4` and updated all module dependencies

## [0.1.11] - 2026-05-25

### Changed

- refreshed `.github/copilot-instructions.md` to add missing `internal/domain/repositories/` directory to the repository structure tree

## [0.1.10] - 2026-05-22

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.9] - 2026-05-19

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.8] - 2026-05-08

### Changed

- changed the Go version to `1.26.3` and updated all module dependencies

## [0.1.7] - 2026-04-28

### Changed

- changed the Go module dependencies to their latest versions
- refreshed `.github/copilot-instructions.md` to remove stale `horusec.json` references, fix constants naming convention, and correct testify description

## [0.1.6] - 2026-04-19

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.5] - 2026-04-15

### Changed

- changed the Go version to `1.26.2` and updated all module dependencies

## [0.1.4] - 2026-04-14

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.3] - 2026-03-24

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.2] - 2026-03-22

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.1] - 2026-03-19

### Changed

- changed the Go module dependencies to their latest versions

## [0.1.0] - 2026-03-12

### Added

- added `https://statusinvest.com.br/etf/eua/{ticker}` website to be crawled because of inaccurate data from `https://dividendhistory.org/`
- added table colors to make it easier to read the data (green for target values, red for opposite)
- added the feature to calculate the dividend yield of an ETF on the last X years
- added the feature to crawl the `https://dividendhistory.org/` website to get the dividend history of an ETF
- added the feature to fetch the `https://api.nasdaq.com/api/quote/{ticker}` endpoint to get the price history of an ETF

### Changed

- changed the Go version to `1.26.1` and updated all module dependencies
- changed the structure to comply DDD standards and principles
- corrected the structure to respect the community standards
