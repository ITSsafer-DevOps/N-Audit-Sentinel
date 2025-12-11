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
```go
// Attach to pod (local operator example)
package main
import (
    "os"
    "os/exec"
)
func main(){
    cmd := exec.Command("kubectl","exec","-it","n-audit-sentinel","--","bash")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    _ = cmd.Run()
}
```
**Expected result:**
- You should see a Kali Linux shell prompt (`┌──(root㉿n-audit-sentinel)-[/var/lib/n-audit]`)
- The `n-audit-sentinel` process may skip TUI prompts if stdin is not attached at startup

### 2. Verify Kali environment
Use the attached shell or programmatic checks. Example (programmatic command run via `kubectl exec`):

```go
// runchecks.go
package main
import (
    "fmt"
    "os/exec"
)
func run(cmd ...string) string { out,_ := exec.Command("kubectl", append([]string{"exec","n-audit-sentinel","--"}, cmd...)...).CombinedOutput(); return string(out) }
func main(){
    fmt.Println(run("whoami"))
    fmt.Println(run("cat","/etc/os-release"))
    fmt.Println(run("uname","-a"))
}
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
To inspect session logs programmatically or from hostPath, see `docs/DEPLOYMENT_HELPERS.md`.

Quick programmatic example (kubectl exec):

```go
out,_ := exec.Command("kubectl","exec","n-audit-sentinel","--","cat","/var/lib/n-audit/session.log").CombinedOutput()
fmt.Println(string(out))
```
**Expected result:**
```
2025-11-30T13:10:19.303705939Z === Infrastructure Discovery ===
2025-11-30T13:10:19.303736389Z K8s API Server: 10.43.0.1:443
2025-11-30T13:10:19.303746967Z DNS Servers: [10.43.0.10]
2025-11-30T13:10:19.303754638Z ================================
```

### 5. Test graceful shutdown (from another terminal)

Programmatic send (exec):

```go
// send-sig.go
package main
import ("os/exec")
func main(){ exec.Command("kubectl","exec","n-audit-sentinel","--","/usr/local/bin/n-audit").Run() }
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
Use the deployment helpers or the Makefile; for automated workflows see `docs/DEPLOYMENT_HELPERS.md` and `TESTING_AND_VERIFICATION.md` for Go-based e2e examples.

Example programmatic recreate + attach (illustrative):

```go
exec.Command("kubectl","delete","pod","n-audit-sentinel","--force","--grace-period=0").Run()
exec.Command("kubectl","create","-f","deploy/k8s/pod.yaml").Run()
exec.Command("kubectl","exec","-it","n-audit-sentinel","--","/usr/local/bin/n-audit-sentinel").Run()
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

To run non-interactive programmatically in Go:

```go
cmd := exec.Command("kubectl","exec","n-audit-sentinel","--","/usr/local/bin/n-audit-sentinel")
cmd.Env = append(os.Environ(), "N_AUDIT_PENTESTER_NAME=John Doe", "N_AUDIT_CLIENT_NAME=Acme Corp")
cmd.Run()
```

---

## Contact
For questions or bug reports:
- Developer: Kristian Kasnik
- Company:
- Website: https://www.example.com
