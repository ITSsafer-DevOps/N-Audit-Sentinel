# N-Audit Sentinel - Manual Testing Guide

## Why manual testing?
The TUI (Text User Interface) interaction runs over a TTY (stdin/stdout), so prompts are **not visible in `kubectl logs`**. The logs contain only the banner and system messages.

## Verified functionality (automated):
✅ **Banner** - ASCII art with magenta "N-Audit Sentinel"  
✅ **Infrastructure Discovery** - K8s API (10.43.0.1:443) and DNS (10.43.0.10)  
✅ **Kali Linux environment** - `/etc/os-release` = Kali GNU/Linux Rolling 2025.4  
✅ **Safety Loop** - the shell restarts after `exit`  
✅ **Graceful Shutdown** - `n-audit` CLI sends SIGUSR1  

---

## Manual test - Full flow

### 1. Attach to the Pod
```bash
kubectl exec -it n-audit-sentinel -- bash
```

**Expected result:**
- You should see a Kali Linux shell prompt (`┌──(root㉿n-audit-sentinel)-[/var/lib/n-audit]`)
- The `n-audit-sentinel` process may skip TUI prompts if stdin is not attached at startup

### 2. Verify Kali environment
```bash
whoami
# Expected output: root

cat /etc/os-release | head -5
# Expected output: Kali GNU/Linux Rolling 2025.4

uname -a
# Expected output: Linux n-audit-sentinel ... (Kali kernel)
```

### 3. Test Safety Loop
```bash
exit
```

**Expected result:**
- The shell restarts (the Pod does not terminate)
- The prompt should reappear: `┌──(root㉿n-audit-sentinel)-[/var/lib/n-audit]`
- The application logs: `[N-Audit] Shell exited with code X. Restarting...`

### 4. Verify session logs
```bash
cat /var/lib/n-audit/session.log
```

**Expected result:**
```
2025-11-30T13:10:19.303705939Z === Infrastructure Discovery ===
2025-11-30T13:10:19.303736389Z K8s API Server: 10.43.0.1:443
2025-11-30T13:10:19.303746967Z DNS Servers: [10.43.0.10]
2025-11-30T13:10:19.303754638Z ================================
```

### 5. Test graceful shutdown (from another terminal)
```bash
# Terminal 1: stay attached to the Pod
kubectl exec -it n-audit-sentinel -- bash

# Terminal 2: send SIGUSR1
kubectl exec n-audit-sentinel -- /usr/local/bin/n-audit
```

**Expected result in Terminal 1:**
- The shell will terminate (the Pod may restart or stop depending on `restartPolicy`)

**Expected logs:**
```bash
kubectl logs n-audit-sentinel --tail=10
# Contains: [N-Audit] Received signal user defined signal 1. Initiating shutdown...
```

---

## Why TUI prompts do not appear in logs

The TUI uses a **TTY (pseudoterminal)** for interactive I/O:
- `stdin: true` and `tty: true` in the Pod spec allow attaching a TTY
- Prompts such as "Pentester Name:" and "Client Name:" are written **only to the TTY**
- `kubectl logs` captures **only stdout/stderr** from the process, not TTY output

### Alternative: Script-based testing
For automated TUI testing you can:
1. Use `expect` or a similar tool
2. Simulate TTY interaction using the `script` command
3. Or implement a non-interactive mode (environment variable `N_AUDIT_NONINTERACTIVE=true`)

---

## Summary of successful tests

### Verified components:
1. ✅ **Docker image build** - Go 1.24 multi-stage build succeeded
2. ✅ **K3s deployment** - Pod is Running and Cilium CNI active
3. ✅ **Infrastructure discovery** - K8s API and DNS detected
4. ✅ **ASCII banner** - displayed in logs with colors
5. ✅ **Kali Linux environment** - shell works and commands are available
6. ✅ **Safety loop** - shell restarts after exit
7. ✅ **Graceful shutdown** - SIGUSR1 correctly handled

### Unverified (requires interactive TTY testing):
- [ ] TUI prompt "Pentester Name:" - display and input handling
- [ ] TUI prompt "Client Name:" - display and input handling
- [ ] TUI prompt "Scope (IP/CIDR)" - validation and double-enter finish
- [ ] TUI prompt "Scope (Domains)" - validation and double-enter finish
- [ ] Network Policy creation - creation of Cilium policy after scope input
- [ ] SSH signature - sign the session log at shutdown (requires SSH key)

---

## Additional steps for full testing

### Option 1: Interactive TTY testing
```bash
# Recreate Pod for manual interaction
kubectl delete pod n-audit-sentinel --force --grace-period=0
kubectl create -f deploy/k8s/pod.yaml
kubectl exec -it n-audit-sentinel -- /usr/local/bin/n-audit-sentinel
```
**Note:** stdin may not be attached at startup; the TUI may skip prompts.

### Option 2: Add debug mode
Edit `cmd/n-audit-sentinel/main.go`:
```go
if os.Getenv("N_AUDIT_DEBUG") == "true" {
    log.Printf("[DEBUG] Waiting 30s for manual attachment...")
    time.Sleep(30 * time.Second)
}
```

### Option 3: Non-interactive testing mode
Set environment variables:
```bash
N_AUDIT_PENTESTER_NAME="John Doe"
N_AUDIT_CLIENT_NAME="Acme Corp"
N_AUDIT_SCOPE_IPS="192.168.1.0/24,10.0.0.0/8"
N_AUDIT_SCOPE_DOMAINS="example.com,test.local"
```

---

## Contact
For questions or bug reports:
- Developer: Kristian Kasnik
- Company:
- Website: https://www.example.com
