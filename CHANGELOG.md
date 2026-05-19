# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-05-19

### Added

- Initial release
- Certificate obtain, renew, revoke, and list commands
- Three interaction modes: novice (wizard), standard (flags), expert (full ACME control)
- HTTP-01 challenge with built-in standalone server
- DNS-01 challenge with Cloudflare, Route53, and DigitalOcean providers
- Automatic DNS provider detection from environment variables
- Secure certificate storage with POSIX permissions (0600 for keys)
- Batch renewal command (`renew-all`) with 30-day expiry threshold
- Auto-renewal scheduling via cron (Linux/macOS) and Task Scheduler (Windows)
- Cross-platform builds for Linux, macOS, and Windows (amd64 and arm64)
- GoReleaser configuration for automated releases
- Unit tests for store, wizard, and DNS provider modules
- Makefile with build, test, release, and clean targets
