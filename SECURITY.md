# Security Policy

This security policy describes responsible disclosure, severity handling, and timelines for the N-Audit Sentinel project.

## Reporting a Vulnerability

If you discover a security vulnerability, please report it privately via email to: itssafer@itssafer.org. Do not open a public issue with vulnerability details.

Include in your report:

- A clear description of the vulnerability
- Steps to reproduce (PoC) if available
- Affected versions and environment details

We will acknowledge receipt within 72 hours and provide a planned remediation timeline.

## Severity and CVSS

We classify vulnerabilities using the Common Vulnerability Scoring System (CVSS) v3.1. Triage will map CVSS scores to priority and SLA:

- Critical (CVSS >= 9.0): Fix within 7 days, coordinated disclosure.
- High (7.0 <= CVSS < 9.0): Fix within 14 days.
- Medium (4.0 <= CVSS < 7.0): Fix within 30 days.
- Low (CVSS < 4.0): Fix within 90 days or as part of scheduled maintenance.

If a vulnerability affects multiple distributions (K3s/Talos/OpenShift), we will coordinate with affected maintainers as appropriate.

## Disclosure Timeline and Embargo

We follow a responsible disclosure process. By default:

- We will acknowledge within 72 hours.
- We will propose a remediation timeline (see SLA above).
- We will coordinate an embargo period to allow fixes to be prepared before public disclosure. Embargo duration will depend on severity and coordination with downstream vendors.

## In-Scope vs Out-of-Scope

In-scope:

- Vulnerabilities in the N-Audit Sentinel source code and official release artifacts
- Packaging scripts and deployment manifests maintained in this repository

Out-of-scope:

- Misconfigurations in downstream platform (e.g. Talos-specific node configuration), unless reproducible solely via N-Audit Sentinel
- Vulnerabilities in third-party dependencies — we will report and track responsibly but will not be the primary maintainer of those fixes

## Remediation and Patching

Fixes will be released through standard GitHub releases and distribution packages. Patches will always include:

- a CVE assignment request if severity warrants
- a security advisory with affected versions and mitigation steps
- checksums for any release artifacts

## Contact & Credits

Contact: itssafer@itssafer.org

We credit reporters who wish to be acknowledged, unless they request anonymity.
# Security Policy

Responsible disclosure

We appreciate reports of security vulnerabilities. To report a vulnerability, please follow these steps:

1. Do not open a public issue. Instead, contact the maintainers privately via email: itssafer@itssafer.org
2. **LinkedIn:** linkedin.com/in/kristián-kašník-03056a377
3. Include a clear description, reproduction steps, affected versions, and any PoC code if available.
4. We will acknowledge receipt within 72 hours and provide a timeline for fixes.

Please do NOT include any sensitive data when reporting publicly.

Disclosure timeline

- We aim to provide a fix and security advisory within 30 days for severities classified as high/critical.
- For coordination and responsible disclosure, we reserve the right to extend the timeline in case of complex fixes.

Credit

We will credit reporters who agree to be acknowledged in the release notes, unless they request anonymity.
