<div align="right">

[🇫🇷 Français](README.fr.md)

</div>

# LeCert

[![CI](https://github.com/valorisa/lecert/actions/workflows/ci.yml/badge.svg)](https://github.com/valorisa/lecert/actions/workflows/ci.yml)
[![Lint](https://github.com/valorisa/lecert/actions/workflows/lint.yml/badge.svg)](https://github.com/valorisa/lecert/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/valorisa/lecert)](https://goreportcard.com/report/github.com/valorisa/lecert)
[![Go Version](https://img.shields.io/github/go-mod/go-version/valorisa/lecert)](https://go.dev/)
[![License](https://img.shields.io/github/license/valorisa/lecert)](https://github.com/valorisa/lecert/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/valorisa/lecert?include_prereleases)](https://github.com/valorisa/lecert/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-brightgreen)](https://github.com/valorisa/lecert#installation)

A cross-platform CLI tool that simplifies Let's Encrypt certificate management for users of all skill levels.
Lecert adapts its interface to three interaction modes (novice, standard, expert) so that both first-time users and seasoned sysadmins can obtain, renew, and revoke TLS certificates without friction.

## What is This and Why Do I Need It?

### The Problem in Plain Language

When you visit a website and see the padlock icon in your browser's address bar, that means the connection between your computer and the website is encrypted. Nobody can spy on what you type (passwords, credit card numbers) or tamper with the page content. This protection is provided by a **TLS certificate** — a small digital file that proves the website is who it claims to be.

If you run your own website, blog, API, or any internet-facing service, you need a TLS certificate. Without one, browsers will show a scary "Not Secure" warning to your visitors, search engines will rank you lower, and sensitive data travels in the open.

### What is Let's Encrypt?

[Let's Encrypt](https://letsencrypt.org/) is a free, nonprofit Certificate Authority that issues TLS certificates at no cost. Before Let's Encrypt existed (2015), certificates cost money and required manual paperwork. Now anyone can get one for free — but the process still involves technical steps that can be intimidating.

### What Does Lecert Do?

Lecert is a command-line tool that handles the entire certificate lifecycle for you:

1. **Obtain** — Proves you own a domain and gets a certificate from Let's Encrypt
2. **Renew** — Automatically refreshes your certificate before it expires (they last 90 days)
3. **Revoke** — Invalidates a certificate if your server is compromised or you lose control of the domain

### Do I Need This?

You need Lecert (or a similar tool) if any of these apply:

- You host a website, app, or API on your own server
- You see "Not Secure" warnings when visiting your site
- Your current certificate expired and you need a new one
- You want HTTPS but don't want to pay for certificates
- You manage multiple domains and want one tool to handle them all

You do NOT need this if:

- You use a hosting provider that handles certificates automatically (Vercel, Netlify, Heroku)
- You already use Caddy (which has built-in automatic HTTPS)
- Your site is purely local/internal with no internet exposure

### How Does Domain Verification Work?

Let's Encrypt needs proof that you actually own the domain before issuing a certificate. This prevents someone from getting a certificate for a domain they don't control. There are two common methods:

**HTTP Challenge (recommended for beginners):** Let's Encrypt asks your server to place a specific file at a specific URL. If your server can respond correctly, it proves you control the domain. This requires port 80 to be open on your server.

**DNS Challenge (for advanced setups):** Instead of placing a file on your server, you add a specific DNS record to your domain's configuration. This is useful when port 80 is blocked, when you use a CDN, or when you need certificates for internal servers that aren't directly reachable from the internet.

### Choosing Your Comfort Level

Lecert offers three modes because people have different experience levels:

| If you are... | Use this mode | What it feels like |
|---------------|---------------|--------------------|
| New to servers and certificates | `--mode novice` | A friendly wizard that asks 3 simple questions |
| Comfortable with command-line tools | `--mode standard` (default) | Familiar flag-based interface like other CLI tools |
| A sysadmin who wants full control | `obtain-expert` command | Every ACME protocol option exposed, nothing hidden |

You can always start with novice mode and graduate to standard or expert as you gain confidence. The certificates produced are identical regardless of which mode you use.

## Table of Contents

- [What is This and Why Do I Need It?](#what-is-this-and-why-do-i-need-it)
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
- [Frequently Asked Questions](#frequently-asked-questions)
- [Glossary](#glossary)
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

## Frequently Asked Questions

### I'm a complete beginner. Will this break my server?

No. Lecert never modifies your web server configuration (Nginx, Apache, etc.). It only obtains certificate files and stores them in a folder. You still need to configure your web server to use those files, but Lecert cannot accidentally break anything already running.

If you use `--staging` mode, the certificates produced are not trusted by browsers but are otherwise identical — this lets you practice the entire flow without any risk or rate limiting.

### What if my certificate expires?

Let's Encrypt certificates are valid for 90 days. Lecert's `schedule install` command sets up automatic renewal that checks daily and renews any certificate expiring within 30 days. If you forget to set up auto-renewal, your site will show a browser warning after 90 days, but nothing is permanently broken — just run `lecert cert renew --domain yourdomain.com` to fix it.

### Do I need to be root/administrator?

Not usually. Lecert stores certificates in your home directory (`~/.lecert/certs/`) and uses an unprivileged HTTP server on port 5002 for challenges. The only case where root might be needed is if you want to use port 80 directly for HTTP challenges (ports below 1024 require elevated privileges on most systems). The recommended approach is to use a reverse proxy or port forwarding instead.

### Can I use this with Nginx/Apache/Caddy?

Yes. Lecert produces standard PEM files (`cert.pem` and `key.pem`). Point your web server configuration at these files:

```nginx
# Nginx example
ssl_certificate     /home/you/.lecert/certs/yourdomain.com/cert.pem;
ssl_certificate_key /home/you/.lecert/certs/yourdomain.com/key.pem;
```

```apache
# Apache example
SSLCertificateFile    /home/you/.lecert/certs/yourdomain.com/cert.pem
SSLCertificateKeyFile /home/you/.lecert/certs/yourdomain.com/key.pem
```

After renewal, reload your web server to pick up the new certificate (`systemctl reload nginx`).

### What's the difference between HTTP and DNS challenges?

| Aspect | HTTP Challenge | DNS Challenge |
|--------|---------------|---------------|
| Requires | Port 80 open on your server | Access to your DNS provider's API |
| Works for | Servers directly reachable from the internet | Any domain, even behind firewalls |
| Difficulty | Easier (just open a port) | Slightly more complex (requires API token setup) |
| Best for | Simple single-server setups | CDN users, internal servers, complex architectures |

If you're not sure, start with HTTP. You can always switch to DNS later.

### I got a rate limit error. What do I do?

Let's Encrypt limits certificate issuance to 5 per domain per week in production. If you're testing or learning, always use `--staging` to avoid hitting these limits. Staging has much higher limits and is designed for experimentation. Once your setup works with staging, remove the `--staging` flag to get a real certificate.

### How is this different from Certbot?

Certbot is the official Let's Encrypt client maintained by the EFF. It's excellent but oriented toward system administrators. Lecert differs in several ways:

- **Single binary** — no Python dependencies, no pip, no virtualenv
- **Novice mode** — a guided wizard that asks 3 questions instead of requiring you to know flags
- **Cross-platform** — same binary works on Linux, macOS, and Windows
- **Does not touch your web server** — Certbot can auto-configure Nginx/Apache, which is convenient but can also break configurations. Lecert only produces files.
- **Lighter scope** — Lecert does one thing (certificate lifecycle) and leaves server configuration to you

Both tools produce the same certificates from the same Let's Encrypt infrastructure.

## Glossary

Terms you might encounter when working with certificates:

| Term | Meaning |
|------|---------|
| **TLS** | Transport Layer Security — the protocol that encrypts web traffic (successor to SSL) |
| **Certificate** | A digital file that proves a server's identity and enables encrypted connections |
| **Private key** | A secret file that only your server should have — never share it |
| **Certificate Authority (CA)** | An organization trusted to issue certificates (Let's Encrypt is a CA) |
| **ACME** | Automatic Certificate Management Environment — the protocol Let's Encrypt uses |
| **Domain** | Your website's address (e.g. `example.com`) |
| **Challenge** | A test that proves you control a domain before a certificate is issued |
| **PEM** | A file format for certificates and keys (the `.pem` files Lecert produces) |
| **Renewal** | Getting a fresh certificate before the current one expires |
| **Revocation** | Invalidating a certificate (e.g. if your server was compromised) |
| **Staging** | A test environment that issues fake certificates for practice |

## Security

If you discover a security vulnerability, please follow the responsible disclosure process described in [SECURITY.md](SECURITY.md). Do not open a public issue for security vulnerabilities.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each release.
