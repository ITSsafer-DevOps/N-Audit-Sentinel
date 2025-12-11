# Verification Guide

Complete testing and validation procedures for N-Audit Sentinel after deployment to Kubernetes.

## Testing Pipeline

```mermaid
flowchart TD
    A["1. Verify Environment<br/>Discovery"] --> B["2. Test Network Policy<br/>Enforcement"]
    B --> C["3. Verify Logging<br/>Integrity"]
    C --> D["4. Check Persistent<br/>Storage"]
    D --> E["5. Graceful Teardown<br/>& Seal Validation"]
    E --> F["✓ All Tests<br/>Passed"]
    
    style A fill:#4A90E2,color:#fff
    style B fill:#4A90E2,color:#fff
    style C fill:#50E3C2,color:#000
    style D fill:#50E3C2,color:#000
    style E fill:#F5A623,color:#000
    style F fill:#7ED321,color:#000
```

## Test 1: Verify Environment Discovery

After the pod starts, confirm that it discovered the Kubernetes API server and DNS resolvers.

### Check pod logs

```bash
kubectl logs n-audit-sentinel
```

**Expected output (in logs):**
```
[N-Audit] Discovered K8s API Server: 10.43.0.1:443
[N-Audit] Discovered DNS Servers: [10.43.0.10]
```

### What's being verified
- Pod can reach the Kubernetes API server
- Pod can discover DNS servers from the cluster
- Logging system is initialized and writing output

**If discovery fails:**
- Pod may not have network access to cluster services
- Check pod's `kubectl logs` for specific errors
- Verify ServiceAccount and RBAC are correctly applied (Test 3)

## Test 2: Verify Network Policy Enforcement (Cilium)

Test that Cilium policies correctly enforce scope boundaries.

### Scenario A: Restricted Scope (Explicit Target)

**Start a session with defined scope:**

```bash
kubectl attach -it n-audit-sentinel -c sentinel
```

**Follow TUI prompts:**
```
Pentester Name: <your-name>
Client Name: <client-name>
[Layer 3] IP/CIDR: 8.8.8.8
(Double Enter)
[Layer 7] Domain: (press Enter twice)
```

**Inside the shell, test policy enforcement:**

```bash
# In-scope should PASS
ping -c 2 8.8.8.8

# Out-of-scope should FAIL
ping -c 2 1.1.1.1

# L7 domain access should FAIL
curl -m 5 -I google.com || echo "curl blocked (expected)"
```

**Expected results:**
- ✓ `ping 8.8.8.8` responds (2/2 packets received)
- ✗ `ping 1.1.1.1` fails (no responses or timeout)
- ✗ `curl google.com` fails or times out (policy blocks)

### Scenario B: Unrestricted Mode (No Scope)

**Start a session without specifying scope:**

```bash
kubectl attach -it n-audit-sentinel -c sentinel
```

**Follow TUI prompts (press Enter through all):**
```
Pentester Name: (press Enter)
Client Name: (press Enter)
[Layer 3] IP/CIDR: (press Enter twice)
[Layer 7] Domain: (press Enter twice)
```

**Check logs for unrestricted mode indicator:**

```bash
kubectl logs n-audit-sentinel --tail=50 | grep -i "unrestricted"
# Expected: message indicating no policy is active
```

**Inside the shell, access should be open:**

```bash
curl -m 5 -I google.com
```

**Expected result:**
- ✓ `curl google.com` succeeds (returns HTTP headers)
- No Cilium policy is enforced

### Debugging Network Policy Issues

```bash
# View applied Cilium policies
kubectl get ciliumnetworkpolicies

# Describe the policy
kubectl describe cnp n-audit-sentinel-policy

# Check Cilium agent logs
kubectl logs -n kube-system -l k8s-app=cilium | grep n-audit
```

## Test 3: Verify Logging Integrity

### Real-time Log Tailing

**From inside the pod:**
```bash
# Execute while session is active
kubectl exec -it n-audit-sentinel -- tail -f /var/lib/n-audit/session.log
```

**From the host (via hostPath):**
```bash
# Execute on the node running the pod
sudo tail -f /mnt/n-audit-data/session.log
```

### Log Format Validation

**Check log structure:**
```bash
kubectl exec n-audit-sentinel -- head -20 /var/lib/n-audit/session.log
```

**Expected log format:**
```
2025-12-11 14:30:45 [Banner] N-Audit Sentinel v1.0.0-Beta
2025-12-11 14:30:45 Pentester: operator
2025-12-11 14:30:45 Client: client-name
2025-12-11 14:30:45 Scope: 8.8.8.8/32
2025-12-11 14:31:00 $ whoami
2025-12-11 14:31:00 root
2025-12-11 14:31:02 $ id
2025-12-11 14:31:02 uid=0(root) gid=0(root) groups=0(root)
```

### Verify ANSI-Free Logging

**Check that logs contain NO escape sequences:**

```bash
# Bad (contains ANSI codes):
# $ grep -E '\x1b\[[0-9;]*m' /var/lib/n-audit/session.log && echo "FAIL: ANSI found"

# Good (no ANSI codes):
kubectl exec n-audit-sentinel -- grep -c $'\x1b' /var/lib/n-audit/session.log
# Expected: 0
```

**All lines must have timestamps:**

```bash
kubectl exec n-audit-sentinel -- \
  awk '!/^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}/' \
  /var/lib/n-audit/session.log | wc -l
# Expected: 0 (no lines without timestamps)
```

## Test 4: Persistent Storage (hostPath)

### Verify Log Persistence

Find the node running the pod and inspect the hostPath:

```bash
# Step 1: Find which node is running the pod
kubectl get pod n-audit-sentinel -o wide
# Output shows NODE column

# Step 2: SSH to that node and inspect hostPath
ssh <node-ip>
sudo ls -lah /mnt/n-audit-data/
sudo tail -n 50 /mnt/n-audit-data/session.log
```

**Expected:**
- `/mnt/n-audit-data/` exists and is readable
- `session.log` mirrors the pod's `/var/lib/n-audit/session.log`
- Log content matches real-time operations

### Verify Log Survives Pod Deletion

```bash
# Note the log content
sudo tail -5 /mnt/n-audit-data/session.log > /tmp/session-before.txt

# Delete the pod
kubectl delete pod n-audit-sentinel

# Log still exists on host
sudo tail -5 /mnt/n-audit-data/session.log > /tmp/session-after.txt

# Verify content is identical
diff /tmp/session-before.txt /tmp/session-after.txt
# Expected: no output (files are identical)
```

## Test 5: Graceful Teardown & Cryptographic Seal

### Trigger Graceful Exit

**From another terminal (while session is active):**

```bash
kubectl exec n-audit-sentinel -c sentinel -- /usr/local/bin/n-audit
```

This sends `SIGUSR1` to PID 1, which triggers:
1. Log finalization
2. SHA256 hash computation
3. SSH signature generation
4. FORENSIC SEAL block appended
5. Cilium policy cleanup
6. Pod termination

### Verify FORENSIC SEAL

**Check the final log block:**

```bash
sudo tail -10 /mnt/n-audit-data/session.log
```

**Expected output:**
```
=== FORENSIC SEAL ===
SHA256 Hash: 3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c
SSH Signature (Base64): AAAAC3NzaC1lZDI1NTE5AAAAgOVCblXV2uXY...
=====================
```

### Validate SHA256 Integrity

**Extract content (excluding seal) and compute hash:**

```bash
# Extract all lines before FORENSIC SEAL
awk '/^=== FORENSIC SEAL ===/{exit} {print}' /mnt/n-audit-data/session.log | sha256sum
# Output: <64-char-hex>  -

# Compare with seal value
grep "SHA256 Hash:" /mnt/n-audit-data/session.log
```

**Expected:**
- Computed hash matches the `SHA256 Hash:` value in the seal block

### Optional: Verify SSH Signature

For full cryptographic verification, extract and validate the SSH signature:

```bash
# Extract public key
sudo cat /mnt/n-audit-data/signing/id_ed25519.pub

# Extract signature from seal
SIGNATURE=$(grep "SSH Signature" /mnt/n-audit-data/session.log | cut -d: -f2)

# Verify with ssh-keygen (requires Go helper or custom validation)
# See README.md for signature verification example
```

## Summary: Complete Verification Checklist

- [ ] **Discovery**: K8s API and DNS found in pod logs
- [ ] **Policy Enforcement**: Restricted scope blocks out-of-scope traffic; unrestricted mode allows all
- [ ] **Logging**: All lines have `YYYY-MM-DD HH:MM:SS` timestamps, zero ANSI codes
- [ ] **Persistence**: Log persists on hostPath even after pod deletion
- [ ] **Seal**: FORENSIC SEAL block appended with valid SHA256 hash
- [ ] **Integrity**: Computed SHA256 matches seal value

**If all checks pass, N-Audit Sentinel is functioning correctly.**

## Troubleshooting Tips

| Symptom | Likely Cause | Check |
|---------|--------------|-------|
| No K8s API discovery | Pod network isolation or API unreachable | `kubectl logs n-audit-sentinel` for errors |
| Cilium policy not applied | RBAC missing or Cilium unhealthy | `kubectl get clusterrolebinding` and `kubectl get ciliumnetworkpolicies` |
| ANSI codes in log | Logger not initialized | Check logger initialization in pod logs |
| No seal block | SSH key missing or unreadable | `ls -la /mnt/n-audit-data/signing/id_ed25519` (must be `600`) |
| Hash mismatch | Log modified after seal | Check file permissions and access logs |

## Cross-References

- **[README.md](README.md)** — Architecture and feature overview
- **[DEPLOYMENT.md](DEPLOYMENT.md)** — Complete deployment instructions
- **[SECURITY.md](SECURITY.md)** — Security model and threat analysis
