# Lecert

A cross-platform CLI tool that simplifies Let's Encrypt certificate management for users of all skill levels.
Lecert adapts its interface to three interaction modes (novice, standard, expert) so that both first-time users and seasoned sysadmins can obtain, renew, and revoke TLS certificates without friction.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Modes](#usage-modes)
- [Commands](#commands)
- [DNS Providers](#dns-providers)
- [Auto-Renewal](#auto-renewal)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)
- [Changelog](#changelog)

## Features

- **Three interaction modes** that adapt to user expertise (novice wizard, standard CLI, expert pass-through)
- **Cross-platform** single binary for Linux, macOS, and Windows with zero runtime dependencies
- **HTTP-01 and DNS-01 challenges** with built-in support for Cloudflare, Route53, and DigitalOcean
- **Automatic renewal scheduling** via cron (Linux/macOS) or Task Scheduler (Windows)
- **Secure key storage** with POSIX permissions (0600) on Unix and ACL restrictions on Windows
- **Batch renewal** of all managed certificates approaching expiration (30-day threshold)
- **Certificate inventory** with status display (valid, expiring soon, expired)
- **Let's Encrypt staging support** for testing without hitting production rate limits
- **GoReleaser integration** for reproducible cross-compiled release builds
- **Environment-based DNS provider detection** for zero-config DNS challenges

## Installation

### From Binary Releases

Download the latest binary for your platform from the [Releases](https://github.com/valorisa/lecert/releases) page.

```bash
# Linux (amd64)
curl -LO https://github.com/valorisa/lecert/releases/latest/download/lecert_linux_amd64.tar.gz
tar xzf lecert_linux_amd64.tar.gz
sudo mv lecert /usr/local/bin/

# macOS (Apple Silicon)
curl -LO https://github.com/valorisa/lecert/releases/latest/download/lecert_darwin_arm64.tar.gz
tar xzf lecert_darwin_arm64.tar.gz
sudo mv lecert /usr/local/bin/

# Windows (amd64) — extract the zip and add to PATH
```

### From Source

Requires Go 1.21 or later.

```bash
go install github.com/valorisa/lecert/cmd/lecert@latest
```

### Build from Repository

```bash
git clone https://github.com/valorisa/lecert.git
cd lecert
make build
# Binary available at ./bin/lecert
```

## Quick Start

### For Beginners (Novice Mode)

The wizard asks exactly three questions and handles everything else automatically.

```bash
lecert --mode novice cert obtain
```

```text
=== Let's Encrypt Certificate Wizard ===

1/3 Domain name (e.g. example.com): mysite.com
2/3 Your email (for certificate expiry alerts): me@mysite.com
3/3 How should we verify domain ownership?
    [1] HTTP (needs port 80 open) — recommended for most setups
    [2] DNS  (needs DNS provider API access)
    Choice [1]: 1

Got it! Requesting certificate for mysite.com via http-01 challenge...

Certificate obtained for mysite.com
  Expires: 2026-08-17
  Stored:  /home/user/.lecert/certs/mysite.com
```

### For Standard Users

```bash
lecert cert obtain --domain mysite.com --email me@mysite.com --staging
```

### For Experts

```bash
lecert cert obtain-expert \
  --domain mysite.com \
  --domain www.mysite.com \
  --email me@mysite.com \
  --challenge dns-01 \
  --dns-provider cloudflare \
  --key-type ec384 \
  --preferred-chain "ISRG Root X1" \
  --timeout 5m
```

## Usage Modes

| Mode | Flag | Audience | Behavior |
|------|------|----------|----------|
| Novice | `--mode novice` | First-time users | Interactive wizard, max 3 questions, sensible defaults |
| Standard | `--mode standard` | Regular users | Flag-based, clear error messages, required flags enforced |
| Expert | `obtain-expert` command | Sysadmins | All ACME options exposed, multi-domain SAN, key type selection |

The mode is set globally via the `--mode` flag and applies to all subcommands.

## Commands

### Certificate Operations

| Command | Description |
|---------|-------------|
| `lecert cert obtain` | Obtain a new certificate (respects current mode) |
| `lecert cert obtain-expert` | Obtain with full ACME control (expert only) |
| `lecert cert renew --domain X` | Renew a specific certificate |
| `lecert cert renew-all` | Renew all certificates expiring within 30 days |
| `lecert cert revoke --domain X` | Revoke a certificate (with interactive confirmation) |
| `lecert cert list` | Display all managed certificates with status |

### Scheduling

| Command | Description |
|---------|-------------|
| `lecert schedule install` | Install auto-renewal (cron or Task Scheduler) |
| `lecert schedule uninstall` | Remove auto-renewal schedule |
| `lecert schedule status` | Check if auto-renewal is active |

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--mode` | `standard` | Interaction mode: novice, standard, expert |
| `--version` | | Display version and exit |
| `--help` | | Display help for any command |

## DNS Providers

Lecert supports DNS-01 challenges via the following providers. Authentication is configured via environment variables.

| Provider | Environment Variables |
|----------|---------------------|
| Cloudflare | `CF_DNS_API_TOKEN` or `CF_API_EMAIL` + `CF_API_KEY` |
| AWS Route53 | `AWS_ACCESS_KEY_ID` + `AWS_SECRET_ACCESS_KEY` + `AWS_REGION` |
| DigitalOcean | `DO_AUTH_TOKEN` |

Lecert auto-detects your DNS provider from environment variables. You can also specify it explicitly:

```bash
export CF_DNS_API_TOKEN="your-token-here"
lecert cert obtain --domain mysite.com --email me@mysite.com --challenge dns-01 --dns-provider cloudflare
```

## Auto-Renewal

Lecert can install a scheduled task that automatically renews certificates approaching expiration.

### Install

```bash
# Daily check at 02:30 (default)
lecert schedule install

# Custom interval
lecert schedule install --interval twice-daily
lecert schedule install --interval hourly
```

### How It Works

The scheduler runs `lecert cert renew-all --quiet` at the configured interval. This command iterates over all managed certificates and renews any that expire within 30 days. Failed renewals are logged to stderr but do not halt the process for remaining certificates.

### Platform Details

| Platform | Mechanism | Schedule Location |
|----------|-----------|-------------------|
| Linux | User crontab | `crontab -l` |
| macOS | User crontab | `crontab -l` |
| Windows | Task Scheduler | `schtasks /Query /TN LecertAutoRenew` |

## Configuration

### Storage Location

Certificates and metadata are stored in `~/.lecert/certs/` by default. Override with the `LECERT_DIR` environment variable.

```bash
export LECERT_DIR=/etc/lecert/certs
lecert cert obtain --domain mysite.com --email me@mysite.com
```

### Directory Structure

```text
~/.lecert/certs/
└── mysite.com/
    ├── cert.pem      # Certificate chain (0644)
    ├── key.pem       # Private key (0600)
    └── meta.json     # Metadata (domain, email, challenge, expiry)
```

### Security

Private keys are stored with restrictive permissions (0600 on Unix). On Windows, standard user-only ACLs apply. Private keys never leave the local filesystem and are never logged or transmitted.

## Architecture

```text
lecert/
├── cmd/lecert/              # CLI entry point and command definitions
│   ├── main.go              # Root command, mode flag, version
│   ├── obtain.go            # Standard obtain + novice routing
│   ├── obtain_expert.go     # Expert mode with all ACME flags
│   ├── renew.go             # Single domain renewal
│   ├── renew_all.go         # Batch renewal (J-30 threshold)
│   ├── revoke.go            # Revocation with confirmation
│   ├── list.go              # Certificate inventory display
│   └── schedule.go          # Scheduler install/uninstall/status
├── internal/
│   ├── acme/                # ACME client wrapper around lego
│   │   ├── acme.go          # Obtain, Renew, Revoke operations
│   │   └── dns.go           # DNS provider registry and detection
│   ├── store/               # Secure certificate storage
│   │   └── store.go         # Save, Load, List with permissions
│   ├── wizard/              # Novice mode interactive wizard
│   │   └── wizard.go        # 3-question guided flow
│   └── scheduler/           # Auto-renewal scheduling
│       ├── scheduler.go     # OS-agnostic dispatcher
│       ├── cron.go          # Linux/macOS crontab management
│       └── windows.go       # Windows Task Scheduler management
├── Makefile                 # Build, test, release targets
├── .goreleaser.yaml         # Release configuration
├── go.mod
└── go.sum
```

### Design Decisions

- **lego library**: Provides mature ACME protocol support with 100+ DNS providers, avoiding reinvention
- **cobra CLI framework**: Industry standard for Go CLIs, enables subcommands and auto-completion
- **internal/ packages**: Prevent external imports, enforce API boundaries
- **No CGO**: Enables true cross-compilation without toolchain dependencies
- **Metadata files (meta.json)**: Allow stateless renewal by persisting email and challenge type per domain

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenience targets)
- GoReleaser (optional, for release builds)

### Build

```bash
make build          # Build for current platform
make release        # Cross-compile all platforms
make test           # Run all tests
make clean          # Remove build artifacts
```

### Cross-Compilation Targets

| OS | Architecture | Binary Name |
|----|-------------|-------------|
| Linux | amd64 | `lecert-linux-amd64` |
| Linux | arm64 | `lecert-linux-arm64` |
| macOS | amd64 | `lecert-darwin-amd64` |
| macOS | arm64 | `lecert-darwin-arm64` |
| Windows | amd64 | `lecert-windows-amd64.exe` |

## Testing

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./internal/store/ -v
go test ./internal/wizard/ -v
go test ./internal/acme/ -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Coverage

| Package | Tests | Coverage |
|---------|-------|----------|
| internal/store | 3 tests (Save/Load, List, Permissions) | Core CRUD operations |
| internal/wizard | 4 tests (valid input, DNS choice, empty domain, EOF) | All wizard paths |
| internal/acme | 5 tests (DNS detection, provider registry) | Provider logic |

### Staging Environment

For integration testing without hitting Let's Encrypt production rate limits, use the `--staging` flag on all commands. Staging certificates are not trusted by browsers but exercise the full ACME flow.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, development workflow, and the process for submitting pull requests.

### Quick Contribution Guide

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Security

If you discover a security vulnerability, please follow the responsible disclosure process described in [SECURITY.md](SECURITY.md). Do not open a public issue for security vulnerabilities.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each release.
