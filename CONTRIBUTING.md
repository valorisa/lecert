# Contributing to Lecert

Thank you for considering contributing to Lecert. This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- Your operating system and version
- Go version (`go version`)
- Lecert version (`lecert --version`)
- Steps to reproduce the behavior
- Expected behavior versus actual behavior
- Any relevant log output or error messages

### Suggesting Features

Feature requests are welcome. Please open an issue with the "enhancement" label and include:

- A clear description of the feature
- The motivation and use case
- Any alternative solutions you considered

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you have added code that should be tested, add tests
3. Ensure the test suite passes (`make test`)
4. Make sure your code builds on all platforms (`make release`)
5. Follow the existing code style (run `gofmt`)
6. Write a clear PR description explaining the change

## Development Setup

```bash
git clone https://github.com/valorisa/lecert.git
cd lecert
go mod tidy
make build
make test
```

## Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation only changes
- `test:` Adding or correcting tests
- `refactor:` Code change that neither fixes a bug nor adds a feature
- `chore:` Maintenance tasks

## Testing

All contributions must include appropriate tests. Run the full suite before submitting:

```bash
go test ./... -v -count=1
```

For integration tests against Let's Encrypt staging, use the `--staging` flag.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
