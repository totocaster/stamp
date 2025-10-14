# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Stamp CLI tool
- Support for multiple note types:
  - Default timestamp (YYYY-MM-DD-HHMM)
  - Daily notes (YYYY-MM-DD)
  - Fleeting notes (YYYY-MM-DD-FHHMMSS)
  - Voice transcripts (YYYY-MM-DD-VTHHMMSS)
  - Analog/slipbox notes (YYYY-MM-DD-AN)
  - Monthly reviews (YYYY-MM)
  - Yearly reviews (YYYY)
  - Project notes (PXXXX)
- Smart counter management for analog and project notes
- Persistent counter storage
- Clipboard support for macOS
- Configuration file support (~/.stamp/config.yaml)
- Dual command names: `stamp` and `nid`
- Cross-platform support (macOS, Linux, Windows)
- Automated release pipeline with GoReleaser
- Homebrew tap for easy installation
- Comprehensive test coverage
- Version command with build information

### Infrastructure
- GitHub Actions CI/CD pipeline
- Automated releases with GoReleaser
- Homebrew formula generation
- Multi-platform builds
- Code linting with golangci-lint

## [0.1.0] - TBD

Initial public release. See Unreleased section for features.

[Unreleased]: https://github.com/totocaster/stamp/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/totocaster/stamp/releases/tag/v0.1.0