#!/bin/bash
# Interactive test script for N-Audit Sentinel TUI

echo "=== N-Audit Sentinel Interactive Test ==="
echo ""
echo "Tento test overí:"
echo "1. Zobrazenie ASCII banneru"
echo "2. Prompt pre Pentester Name"
echo "3. Prompt pre Client Name"
echo "4. Prompt pre Scope (IP/CIDR)"
echo "5. Prompt pre Scope (Domains)"
echo "6. Spustenie Kali Linux shellu"
echo "7. Safety Loop (exit = reštart shellu)"
echo ""
echo "Pripájam sa do podu..."
echo ""

kubectl exec -it n-audit-sentinel -- bash -c '
echo "=== Kontrola procesu n-audit-sentinel ==="
ps aux | grep n-audit-sentinel | grep -v grep
echo ""
echo "=== Session log ==="
cat /var/lib/n-audit/session.log 2>/dev/null || echo "Žiadny session.log"
echo ""
echo "=== Teraz otvorím Kali shell - MANUÁLNE TESTOVANIE: ==="
echo "1. Zadaj \"whoami\" - overíš Kali environment"
echo "2. Zadaj \"cat /etc/os-release | head -5\" - overíš Kali Linux"
echo "3. Zadaj \"exit\" - overíš safety loop (shell by sa mal reštartovať)"
echo ""
exec bash
'
