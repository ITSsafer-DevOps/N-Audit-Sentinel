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

## Threat Model (High Level)

```mermaid
graph TB
	User["👤 Pentester / Auditor"] -->|kubectl attach| Pod["🔷 N-Audit Sentinel Pod (PID 1)"]
	Pod -->|writes| Host["💾 HostPath (/mnt/n-audit-data)"]
	Pod -->|applies| Cilium["🔌 Cilium CNI (CiliumNetworkPolicy)"]
	Attacker["⚠️ Malicious Container"] -. attempt -> Pod
	Attacker -. lateral -> Host
	Note["Threats: Key theft, log tampering, RBAC misuse, eBPF escalation"]
	style Pod fill:#4A90E2,stroke:#2E5C8A,color:#fff
	style Cilium fill:#F5A623,stroke:#B8770B,color:#000
	style Host fill:#50E3C2,stroke:#2E8B74,color:#000
	style Attacker fill:#D0021B,stroke:#8B0000,color:#fff
```

## Operational Security Checks (Go)

Example: verify signing key permissions and presence programmatically.

```go
package main
import (
	"fmt"
	"os"
)
func main(){
	info, err := os.Stat("/mnt/n-audit-data/signing/id_ed25519")
	if err != nil { fmt.Println("missing key") ; os.Exit(2) }
	if info.Mode().Perm() != 0600 { fmt.Println("bad perms", info.Mode()) ; os.Exit(3) }
	fmt.Println("private key present and permissions OK")
}
```
