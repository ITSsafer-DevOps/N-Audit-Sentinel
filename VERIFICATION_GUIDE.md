# N-Audit Sentinel – Verification Guide

This guide verifies functionality and integrity on a running cluster.

```mermaid
flowchart TB
	A[Attach to pod] --> B[Fill Pentester/Client]
	B --> C[Define Scope (IPs/Domains)]
	C --> D[Test in-scope passes]
	C --> E[Test out-of-scope blocked]
	D & E --> F[Tail session.log]
	F --> G[n-audit exit]
	G --> H[Seal block appended]
	H --> I[Verify SHA-256 & signature]
```

<!-- mermaid: validated -->

## 1. Verifying Environment Discovery
- After the pod starts, check logs:
```bash
kubectl logs n-audit-sentinel
```
- Confirm these lines appear:
- `[N-Audit] Discovered K8s API Server: <IP>:<PORT>`
- `[N-Audit] Discovered DNS Servers: [<dns-ip-1>, <dns-ip-2>, ...]`

If any are missing, DNS or API access may be misconfigured.

## 2. Verifying Network Policy Enforcement (Cilium)

### Scenario A: Restricted Scope
- Start a session with defined scope. Attach to the container and provide scope (example target IP `8.8.8.8`):
```bash
kubectl attach -it n-audit-sentinel -c sentinel
# Fill prompts:
# Pentester Name: <your name>
# Client Name: <client>
# [Layer 3] IP/CIDR: 8.8.8.8
# (Double Enter to continue)
# [Layer 7] Domain: (press Enter twice to finish)
```
- Inside the shell, test policy enforcement:
```bash
# In-scope should PASS
ping -c 2 8.8.8.8

# Out-of-scope should FAIL
ping -c 2 1.1.1.1

# L7 domain access should FAIL
curl -m 5 -I google.com || echo "curl blocked (expected)"
```
Expected:
- Ping `8.8.8.8` works.
- Ping `1.1.1.1` fails (no replies/timeouts).
- `curl google.com` blocked.

### Scenario B: Unrestricted Mode
- Start session without scope (press Enter through all prompts):
```bash
kubectl attach -it n-audit-sentinel -c sentinel
# Pentester Name: <enter>
# Client Name: <enter>
# [Layer 3] (press Enter twice)
# [Layer 7] (press Enter twice)
```
- Confirm log indicates unrestricted mode:
```bash
kubectl logs n-audit-sentinel --tail=50 | grep "Running in unrestricted mode"
```
- Inside shell, access should be open:
```bash
curl -m 5 -I google.com
```
Expected:
- `curl google.com` works; no Cilium policy is applied.

## 3. Verifying Logging Integrity
- View the session log in real-time:
```bash
kubectl exec -it n-audit-sentinel -- tail -f /var/lib/n-audit/session.log
```
Check for:
- No ANSI color codes (e.g., `\x1b[35m`).
- Every line starts with `YYYY-MM-DD HH:MM:SS` (UTC), e.g., `2025-11-30 15:10:51`.

If color codes appear, log routing may be incorrect; only clean text must be present.

## 4. Host-Level Log Persistence (hostPath)
When using a hostPath mount, the log persists on the node even after pod deletion.

Steps:
1. Find the node running the pod:
```bash
kubectl get pod n-audit-sentinel -o wide
```
2. On that node, inspect the mount:
```bash
sudo ls -lah /mnt/n-audit-data/
sudo tail -n 50 /mnt/n-audit-data/session.log
```
Expected:
- `session.log` exists and mirrors `/var/lib/n-audit/session.log` inside the pod.

## 5. Graceful Seal and Validation
End the session from another terminal:
```bash
kubectl exec n-audit-sentinel -- /usr/local/bin/n-audit
```
The log appends:
```
=== FORENSIC SEAL ===
SHA256 Hash: <hex>
SSH Signature (Base64): <base64>
=====================
```

Validate SHA‑256 for the session content (exclude the seal block):
```bash
awk '/^=== FORENSIC SEAL ===/{exit} {print}' /mnt/n-audit-data/session.log | sha256sum
```
Compare the computed hex with the `SHA256 Hash:` line.

Optionally verify the SSH signature using the public key via a short Go helper (example provided in `README.md`).

## Cross‑Links
- Deployment: `DEPLOYMENT.md`
- Overview and security model: `README.md`

## How to verify a release artifact (SHA256)

1. Download the release tarball and the `.sha256` file from the GitHub Release page.

```bash
# example
curl -L -o n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz \
	https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases/download/v1.0.0-Beta/n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz
curl -L -o n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256 \
	https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/releases/download/v1.0.0-Beta/n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256
```

2. Compute local SHA256 and compare:

```bash
sha256sum n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz
cat n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz.sha256
# Both values should match (hex string)
```

3. (Optional) Extract and run a smoke test:

```bash
tar -xzf n-audit-sentinel-v1.0.0-Beta-linux-amd64.tar.gz
./n-audit-sentinel --system-audit || echo "smoke test failed"
```

If the SHA256 matches and the smoke test succeeds, the artifact is verified.
