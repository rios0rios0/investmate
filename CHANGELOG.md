# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

When a new release is proposed:

1. Create a new branch `bump/x.x.x` (this isn't a long-lived branch!!!);
2. The Unreleased section on `CHANGELOG.md` gets a version number and date;
3. Open a Pull Request with the bump version changes targeting the `main` branch;
4. When the Pull Request is merged, a new git tag must be created using [GitHub environment](https://github.com/rios0rios0/pipelines/tags).

Releases to productive environments should run from a tagged version.
Exceptions are acceptable depending on the circumstances (critical bug fixes that can be cherry-picked, etc.).

## [Unreleased]

### Added

- added the feature to crawl the `https://dividendhistory.org/` website to get the dividend history of an ETF
- added the feature to fetch the `https://api.nasdaq.com/api/quote/{ticker}` endpoint to get the price history of an ETF
- added the feature to calculate the dividend yield of an ETF on the last X years
- added table colors to make it easier to read the data (green for target values, red for opposite)
- added `https://statusinvest.com.br/etf/eua/{ticker}` website to be crawled because of inaccurate data from `https://dividendhistory.org/`

### Changed

- corrected the structure to respect the community standards
- changed the structure to comply DDD standards and principles
