# N-Audit Sentinel - Manual Testing Guide

## Prečo manuálne testovanie?
TUI (Text User Interface) interakcia prebieha cez TTY (stdin/stdout), preto sa **nezobrazuje v kubectl logs**. Logy obsahujú len banner a systémové správy.

## Overené funkcionality (automaticky):
✅ **Banner** - ASCII art s magenta "N-Audit Sentinel"  
✅ **Infrastructure Discovery** - K8s API (10.43.0.1:443) a DNS (10.43.0.10)  
✅ **Kali Linux environment** - `/etc/os-release` = Kali GNU/Linux Rolling 2025.4  
✅ **Safety Loop** - shell sa reštartuje po `exit`  
✅ **Graceful Shutdown** - `n-audit` CLI posiela SIGUSR1  

---

## Manuálny test - Kompletný flow

### 1. Pripojenie do podu
```bash
kubectl exec -it n-audit-sentinel -- bash
```

**Očakávaný výsledok:**
- Otvorí sa Kali Linux shell (`┌──(root㉿n-audit-sentinel)-[/var/lib/n-audit]`)
- Proces `n-audit-sentinel` už preskočil TUI (pretože stdin nie je pripojený pri štarte)

### 2. Overenie Kali environmentu
```bash
whoami
# Výstup: root

cat /etc/os-release | head -5
# Výstup: Kali GNU/Linux Rolling 2025.4

uname -a
# Výstup: Linux n-audit-sentinel ... (Kali kernel)
```

### 3. Test Safety Loop
```bash
exit
```

**Očakávaný výsledok:**
- Shell sa **reštartuje** (neukončí sa pod)
- Znova sa zobrazí prompt: `┌──(root㉿n-audit-sentinel)-[/var/lib/n-audit]`
- Aplikácia loguje: `[N-Audit] Shell exited with code X. Restarting...`

### 4. Overenie session logov
```bash
cat /var/lib/n-audit/session.log
```

**Očakávaný výsledok:**
```
2025-11-30T13:10:19.303705939Z === Infrastructure Discovery ===
2025-11-30T13:10:19.303736389Z K8s API Server: 10.43.0.1:443
2025-11-30T13:10:19.303746967Z DNS Servers: [10.43.0.10]
2025-11-30T13:10:19.303754638Z ================================
```

### 5. Test graceful shutdown (z iného terminálu)
```bash
# Terminál 1: zostane v pode
kubectl exec -it n-audit-sentinel -- bash

# Terminál 2: pošle SIGUSR1
kubectl exec n-audit-sentinel -- /usr/local/bin/n-audit
```

**Očakávaný výsledek v Termináli 1:**
- Shell sa ukončí (pod sa reštartuje alebo stopne, závisí od `restartPolicy`)

**Očakávaný výsledek v logoch:**
```bash
kubectl logs n-audit-sentinel --tail=10
# Obsahuje: [N-Audit] Received signal user defined signal 1. Initiating shutdown...
```

---

## Prečo TUI prompty nie sú v logoch?

TUI používa **TTY (pseudoterminal)** pre interaktívnu komunikáciu:
- `stdin: true` a `tty: true` v Pod spec umožňujú pripojiť TTY
- Prompty ako "Pentester Name:", "Client Name:" sa vypisujú **len do TTY**
- `kubectl logs` zachytáva **len stdout/stderr** z procesu, nie z TTY

### Alternatíva: Script-based testing
Pre automatizované testovanie TUI by bolo potrebné:
1. Použiť `expect` alebo podobný nástroj
2. Simulovať TTY interakciu cez `script` príkaz
3. Alebo implementovať non-interactive režim (env variable `N_AUDIT_NONINTERACTIVE=true`)

---

## Zhrnutie úspešného testovania

### Overené komponenty:
1. ✅ **Docker image build** - Go 1.24, multi-stage build úspešný
2. ✅ **K3s deployment** - pod Running, Cilium CNI aktívna
3. ✅ **Infrastructure discovery** - K8s API a DNS detekované
4. ✅ **ASCII banner** - zobrazený v logoch s farbami
5. ✅ **Kali Linux environment** - shell funguje, príkazy dostupné
6. ✅ **Safety loop** - shell sa reštartuje po exit
7. ✅ **Graceful shutdown** - SIGUSR1 správne zachytený

### Neoverené (vyžaduje manuálne testovanie s TTY):
- [ ] TUI prompt "Pentester Name:" - zobrazenie a čítanie vstupu
- [ ] TUI prompt "Client Name:" - zobrazenie a čítanie vstupu
- [ ] TUI prompt "Scope (IP/CIDR)" - validácia a double-enter ukončenie
- [ ] TUI prompt "Scope (Domains)" - validácia a double-enter ukončenie
- [ ] Network Policy creation - vytvorenie Cilium policy po zadaní scope
- [ ] SSH signature - podpis session logu pri shutdowne (vyžaduje SSH kľúč)

---

## Ďalšie kroky pre plné testovanie

### Možnosť 1: Interaktívne TTY testovanie
```bash
# Vytvoriť nový pod s manual interakciou
kubectl delete pod n-audit-sentinel --force --grace-period=0
kubectl create -f deploy/k8s/pod.yaml
kubectl exec -it n-audit-sentinel -- /usr/local/bin/n-audit-sentinel
```
**Problém:** stdin nie je pripojený pri štarte, TUI okamžite preskočí prompty.

### Možnosť 2: Pridať debug režim
Upraviť `cmd/n-audit-sentinel/main.go`:
```go
if os.Getenv("N_AUDIT_DEBUG") == "true" {
    log.Printf("[DEBUG] Waiting 30s for manual attachment...")
    time.Sleep(30 * time.Second)
}
```

### Možnosť 3: Non-interactive testing mode
Implementovať env variables:
```bash
N_AUDIT_PENTESTER_NAME="John Doe"
N_AUDIT_CLIENT_NAME="Acme Corp"
N_AUDIT_SCOPE_IPS="192.168.1.0/24,10.0.0.0/8"
N_AUDIT_SCOPE_DOMAINS="example.com,test.local"
```

---

## Kontakt
Pre otázky alebo bug reporty:
- Developer: Kristián Kašník
- Company: 
- Website: https://www..com
