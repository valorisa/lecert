# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest | Yes |
| < latest | No |

## Reporting a Vulnerability

If you discover a security vulnerability in Lecert, please report it responsibly.

### Process

1. **Do NOT open a public issue** for security vulnerabilities
2. Send a description of the vulnerability to the repository owner via GitHub private vulnerability reporting
3. Include steps to reproduce the issue
4. Include the potential impact assessment

### What to Expect

- Acknowledgment of your report within 48 hours
- A plan for addressing the vulnerability within 7 days
- Credit in the security advisory (unless you prefer anonymity)

### Scope

The following are in scope for security reports:

- Private key exposure or leakage
- ACME account key compromise
- Privilege escalation via the scheduler
- Command injection via user inputs
- Insecure file permissions on certificate storage
- Man-in-the-middle vulnerabilities in ACME communication

### Out of Scope

- Vulnerabilities in Let's Encrypt infrastructure
- Vulnerabilities in upstream dependencies (report those to the respective projects)
- Social engineering attacks
- Denial of service via rate limiting (this is a Let's Encrypt feature)
