**Origin Note:** This project was architected and developed as a **proactive R&D initiative** ("going the extra mile") within the recruitment process.

# N-Audit Sentinel

```
███╗   ██╗       █████╗ ██╗   ██╗██████╗ ██╗████████╗    ███████╗███████╗███╗   ██╗████████╗██╗███╗   ██╗███████╗██╗     
████╗  ██║      ██╔══██╗██║   ██║██╔══██╗██║╚══██╔══╝    ██╔════╝██╔════╝████╗  ██║╚══██╔══╝██║████╗  ██║██╔════╝██║     
██╔██╗ ██║█████╗███████║██║   ██║██║  ██║██║   ██║       ███████╗█████╗  ██╔██╗ ██║   ██║   ██║██╔██╗ ██║█████╗  ██║     
██║╚██╗██║╚════╝██╔══██║██║   ██║██║  ██║██║   ██║       ╚════██║██╔══╝  ██║╚██╗██║   ██║   ██║██║╚██╗██║██╔══╝  ██║     
██║ ╚████║      ██║  ██║╚██████╔╝██████╔╝██║   ██║       ███████║███████╗██║ ╚████║   ██║   ██║██║ ╚████║███████╗███████╗
╚═╝  ╚═══╝      ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚═╝   ╚═╝       ╚══════╝╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝
```

**Developer:** Kristián Kašník  
**Contact:** itssafer@itssafer.org  
**License:** MIT License (Open Source)

### Repository Contents
This repository contains:
- Complete Go source code for the N-Audit Sentinel application
- Internal packages with comprehensive test coverage (49% overall, 91%+ signature/policy/tui)
- Go-based utilities for deterministic builds and packaging (`release-manager`, `backup-manager`)
- Deployment manifests and Terraform configurations
- Comprehensive documentation and verification guides

**Note:** Binary release artifacts (`.tar.gz`) are produced by the release pipeline and stored as GitHub Releases, not in Git history. Legacy shell scripts have been removed; use the Go utilities under `cmd/` for all build, packaging, and archival operations.

N-Audit Sentinel is a Kubernetes‑native forensic wrapper that runs as PID 1 inside a Kali Linux pod. It hardens network access with Cilium, guarantees clean and human‑readable logs, and seals every session with a cryptographic signature.

## Core Principles

**What:** N-Audit Sentinel is a Kubernetes-native forensic wrapper that provides controlled network access, cryptographic audit trails, and tamper-evident session recording.

**Why:** Penetration testers need strong guarantees that their session scope is enforced, all commands are logged cleanly, and logs cannot be modified undetected. Cilium policies provide granular network enforcement; cryptographic seals provide integrity verification.

**How:** The application runs as PID 1 in a Kali Linux pod, discovers the Kubernetes API and DNS at startup, accepts scope from a TUI, applies Cilium network policies, and records all PTY activity with timestamps and a final cryptographic seal.

## Go-Based Release & Backup Utilities

The project includes two production-grade Go utilities that replace legacy shell scripts:

### release-manager
**Location:** `cmd/release-manager`  
**Purpose:** Builds binaries with deterministic flags, packages them into `.tar.gz` archives, and generates SHA256 checksums.

**Example:**
```bash
go run ./cmd/release-manager --version v1.0.0-Beta --out out
# Creates: n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz
#          n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256
```

### backup-manager
**Location:** `cmd/backup-manager`  
**Purpose:** Creates a Gold Master source archive using `git archive` and generates SHA256 checksums for reproducibility.

**Example:**
```bash
go run ./cmd/backup-manager --out gold-master-20251210T235959Z.tar.gz --ref HEAD
```

Both utilities are tested, minimal, and designed for local development and CI pipelines. See `cmd/*/README.md` for complete usage documentation.
## Features at a Glance

| Capability | What it does | Why it matters |
|---|---|---|
| PID 1 Safety Loop | Respawns the shell on `exit`/Ctrl+D | Prevents accidental session loss |
| 3‑Zone Cilium Policy | Infra (API+DNS), Maintenance (HTTP/S to repos), Target (your scope) | Blocks out‑of‑scope traffic by default |
| Scope‑Driven Enforcement | IP/CIDR + Domains via TUI (double‑enter) | Clear, explicit authorization surface |
| Clean Logging | Strips ANSI; timestamps each line as `YYYY‑MM‑DD HH:MM:SS` | Human‑readable, diff‑friendly logs |
| Real‑time Flush | `O_SYNC` file writes | Live tailing via SSHFS without delays |
| Cryptographic Seal | Appends SHA‑256 and SSH signature at teardown | Tamper‑evident audit trail |
| Kubernetes‑Aware | Auto‑discovers API and DNS | No manual wiring for cluster basics |
| Open Source | MIT licensed | Transparent and extensible |

## Architecture

### Kubernetes Integration

```mermaid
flowchart LR
  subgraph K8sCluster [\"Kubernetes Cluster\"]
    direction TB
    Pod([\"N-Audit Sentinel Pod<br/>(PID 1 Runtime)\"])
    Vol[(\"hostPath Volume<br/>/var/lib/n-audit\")]
    API[(\"Kubernetes API<br/>Service Discovery\")]
    DNS[(\"CoreDNS<br/>Domain Resolution\")]
    CNP[[\"Cilium Network Policy<br/>Scope Enforcement\"]]
    
    Pod -->|auto-discover| API
    Pod -->|auto-discover| DNS
    Pod -->|create/delete| CNP
    Pod <-->|mount| Vol
  end

  User[\"Pentester/Auditor\"] -->|kubectl attach TTY| Pod
  
  style Pod fill:#4A90E2,color:#fff,stroke:#2E5C8A,stroke-width:2px
  style Vol fill:#50E3C2,color:#000,stroke:#2E8B74,stroke-width:2px
  style CNP fill:#F5A623,color:#000,stroke:#B8770B,stroke-width:2px
```

### Core Modules

| Module | Location | Purpose |
|--------|----------|---------|
| **PID 1 Runtime** | `cmd/n-audit-sentinel` | Process lifecycle, signal handling, session teardown |
| **Exit Helper** | `cmd/n-audit` | Triggers graceful shutdown via SIGUSR1 |
| **Logger** | `internal/logger` | ANSI stripping, per-line timestamp injection (`YYYY-MM-DD HH:MM:SS UTC`) |
| **Policy Engine** | `internal/policy` | Cilium policy generation, validation, and enforcement |
| **PTY Recorder** | `internal/recorder` | Terminal session capture with safety loop on exit |
| **TUI** | `internal/tui` | Interactive banner and scope prompts (double-Enter to finalize) |
| **Validation** | `internal/validation` | IP/CIDR/domain normalization and guardrails |
| **Signature** | `internal/signature` | SHA256 hashing and SSH cryptographic sealing |
| **Discovery** | `internal/discovery` | Kubernetes API and DNS auto-discovery |
| **Backup Manager** | `cmd/backup-manager` | Gold Master archival with checksums |
| **Release Manager** | `cmd/release-manager` | Deterministic artifact packaging |

### Session Lifecycle (Happy Path)

```mermaid
sequenceDiagram
    autonumber
    actor User as Operator
    participant Sentinel as "PID 1<br/>(Sentinel)"
    participant K8s as "K8s API"
    participant DNS as "CoreDNS"
    participant Cilium as "Cilium CNI"
    participant Bash as "/bin/bash<br/>(PTY)"

    User->>Sentinel: kubectl run pod
    Sentinel->>K8s: Auto-discover API IP:PORT
    Sentinel->>DNS: Auto-discover DNS servers
    Sentinel->>Sentinel: Initialize logger (strip ANSI)
    Sentinel->>User: Display banner + prompts
    User->>Sentinel: Enter Pentester name
    User->>Sentinel: Enter Client name
    User->>Sentinel: Enter scope (IPs + Domains)
    Sentinel->>Cilium: Create 3-zone policy<br/>(Infra, Maintenance, Target)
    Sentinel->>Bash: Spawn shell with safety loop
    User->>Bash: Execute audit commands
    Bash->>Bash: Log commands (no ANSI, timestamps)
    User->>Sentinel: n-audit exit (SIGUSR1)
    Sentinel->>Sentinel: Compute SHA256 of log
    Sentinel->>Sentinel: SSH-sign hash
    Sentinel->>Sentinel: Append FORENSIC SEAL
    Sentinel->>Cilium: Delete 3-zone policy
    Sentinel->>User: Exit PID 1 (pod terminates)
```

## Quick Start

For complete deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md). For verification and testing, see [VERIFICATION_GUIDE.md](VERIFICATION_GUIDE.md).

### Minimal 5-Step Workflow

1. **Build Release**
   ```bash
   make release VERSION=v1.0.0-Beta
   ```

2. **Prepare Node Storage**
   ```bash
   sudo mkdir -p /mnt/n-audit-data/signing
   sudo ssh-keygen -t ed25519 -N "" -f /mnt/n-audit-data/signing/id_ed25519 -C "n-audit"
   sudo chmod 600 /mnt/n-audit-data/signing/id_ed25519
   ```

3. **Create ServiceAccount + RBAC**
   ```bash
   kubectl apply -f beta-test-deployment/serviceaccount.yaml
   ```

4. **Deploy Pod**
   ```bash
   kubectl apply -f beta-test-deployment/pod-fixed.yaml
   ```

5. **Attach and Operate**
   ```bash
   kubectl attach -it n-audit-sentinel -c sentinel
   # Follow TUI prompts: Pentester, Client, Scope (IPs/Domains)
   # Exit: n-audit exit (graceful teardown with seal)
   ```

## Runtime Behavior
- TUI captures Pentester and Client and your explicit scope.
- Logs are written to `/var/lib/n-audit/session.log` inside the pod (typically host‑mounted at `/mnt/n-audit-data/session.log`).
- Safety loop: typing `exit` or Ctrl+D respawns the shell. To end, run `n-audit exit` (sends SIGUSR1).
- Teardown appends a “FORENSIC SEAL” block with SHA‑256 and an SSH signature.

## Security Considerations

Why these choices
- Why PID 1: To control lifecycle and signals reliably; avoid accidental termination; guarantee teardown sealing.
- Why root by default (configurable): Many forensic tasks require elevated capabilities. If not needed, run as a non‑root user and review resulting limitations.
- Why Cilium: Granular, identity‑aware network policies and high‑quality L3/L7 enforcement for scope‑driven lockdown.

Forensic guarantees
- Human‑readable timestamps: each line prefixed with `YYYY‑MM‑DD HH:MM:SS` (UTC).
- ANSI‑free logs: all terminal escape sequences are removed before persistence.
- Real‑time writes: `O_SYNC` flush enables live tailing from mounted storage.
- Cryptographic sealing: SHA‑256 of the session content + OpenSSH signature appended at teardown.

Operational notes
- Private key path: set `SSH_SIGN_KEY_PATH` (e.g., `/var/lib/n-audit/signing/id_ed25519`), permissions `700` dir and `600` key.
- ServiceAccount + RBAC: required to create/delete `ciliumnetworkpolicies.cilium.io`.
- Label selectors: the pod must be labeled `app: n-audit-sentinel` to match the policy.

## Build & Release

### Local Development Build
```bash
# Build individual binaries
go build -o bin/n-audit-sentinel ./cmd/n-audit-sentinel
go build -o bin/n-audit ./cmd/n-audit

# Run tests
go test ./...
```

### Release Build (Production)
```bash
# Create release artifacts with Makefile
make release VERSION=v1.0.0-Beta
```

**Output Artifacts:**
- `n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz` — Binary archive
- `n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256` — SHA256 checksum
- `gold-master-<timestamp>.tar.gz` — Source code archive (optional)

### Verification
```bash
# Verify SHA256
sha256sum -c n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256

# Extract and test
tar -xzf n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz
./n-audit-sentinel --version
```

## Deployment Methods

### Method 1: Manifest-Based (kubectl)
Manifest files are provided in `beta-test-deployment/`:
- `serviceaccount.yaml` — ServiceAccount and RBAC bindings for Cilium
- `pod-fixed.yaml` — Complete pod definition with mounts and env vars

```bash
kubectl apply -f beta-test-deployment/serviceaccount.yaml
kubectl apply -f beta-test-deployment/pod-fixed.yaml
```

### Method 2: Terraform (Reproducible)
For automated, version-controlled deployments, use `deploy/terraform/`:

```bash
cd deploy/terraform
terraform init
terraform apply -auto-approve \
  -var="namespace=default" \
  -var="image_name=n-audit-sentinel" \
  -var="image_tag=v1.0.0-Beta"
```

Variables support customization of image, storage class, and namespace. See `deploy/terraform/variables.tf` for details.

## License

MIT License © Kristián Kašník - ITSsafer-DevOps and contributors.  
See [LICENSE](LICENSE) for full text.

---

## Additional Resources

- **[DEPLOYMENT.md](DEPLOYMENT.md)** — Full deployment and configuration guide
- **[VERIFICATION_GUIDE.md](VERIFICATION_GUIDE.md)** — Testing and validation procedures
- **[SECURITY.md](SECURITY.md)** — Security model and operational guidelines
- **[docs/TOOLS.md](docs/TOOLS.md)** — Go utility reference (release-manager, backup-manager)

## Local systemd Deployment (Optional)

For running N-Audit Sentinel as a local systemd service (without Kubernetes), install the provided unit file:

```bash
sudo cp deploy/n-audit-sentinel.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now n-audit-sentinel.service

# Monitor logs
sudo journalctl -u n-audit-sentinel -f
```

### Configuration

**Runtime Parameters:**
- `--system-audit` — Perform basic OS audit scan
- `--verbose` — Enable verbose logging for troubleshooting
- `--log-file /var/lib/n-audit/session.log` — Explicit log file location
- `--sign-key /var/lib/n-audit/signing/id_ed25519` — Path to SSH signing key

**Key Requirements:**
- Signing key must exist and be protected: `chmod 600`
- Signing directory must be restricted: `chmod 700`
- Run as `root` for full capabilities (or review required capabilities for non-root)
- See [deploy/n-audit-sentinel.service](deploy/n-audit-sentinel.service) for complete unit configuration
