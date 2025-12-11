# N-Audit Sentinel — Audit Action Items (Zoznam na spustenie)

## 🗑️ SÚBORY NA VYMAZANIE (Redundancia)

```
Priority  Súbor                                  Veľkosť  Dôvod
─────────────────────────────────────────────────────────────────────────
🔴 HIGH   FINAL_PROJECT_STATUS_v2.md            6.7 KB   Duplikát v2 (zúžená)
🔴 HIGH   ARCHITECTURE_SUPPORT_MATRIX.md        1.5 KB   Duplikuje ENTERPRISE_LEVEL_AUDIT.md
🟠 MED    TEST_SUITE_VERIFICATION.md            1.7 KB   Duplikuje PROJECT_AUDIT_REPORT.md
🟠 MED    COVERAGE_REPORT.md                    658 B    Zastarané metriky (49.5%)
🟡 LOW    FINALIZATION_REPORT.md                2.5 KB   Redundantný s FINAL_PROJECT_STATUS.md
🟡 LOW    RELEASE_ARTIFACTS.md                  270 B    Informatívna len (zvážiť archivovať)
```

**Príkazy na zmazanie:**
```bash
rm -f FINAL_PROJECT_STATUS_v2.md \
      ARCHITECTURE_SUPPORT_MATRIX.md \
      TEST_SUITE_VERIFICATION.md \
      COVERAGE_REPORT.md
```

---

## 🔤 SÚBORY NA LOKALIZÁCIU (Oprava na angličtinu)

```
Súbor                                      Riadky   Zmeny
──────────────────────────────────────────────────────────────────
COMPREHENSIVE_ENHANCEMENT_REPORT.md        84, 127  2 zmeny
```

### Detaily:

**Súbor:** `COMPREHENSIVE_ENHANCEMENT_REPORT.md`

**Zmena 1 (Riadok 84):**
```diff
- 3-vrstvová architektura
+ 3-layer architecture
```

**Zmena 2 (Riadok 127):**
```diff
- BEZ OTÁZOK. AUTONOMNE. 100% HOTOVO.
+ NO QUESTIONS. AUTONOMOUS. 100% COMPLETE.
```

---

## ⚙️ MAKEFILE PRÍKAZY NA PREPIS (Go-native)

### Príkaz 1: Dependency Verification (Riadky 56-60)
**Aktuálne:** Bash `command -v` checks
**Go-native:** `exec.LookPath()` v novom helper CLI

```bash
# Vytvorit:
cmd/verify-deps/main.go

# Potom Makefile:
verify-deps:
    @echo "Verifying dependencies: go, docker, kubectl"
    go run ./cmd/verify-deps
```

---

### Príkaz 2: Linting (Riadok 31)
**Aktuálne:** `if command -v golangci-lint`
**Go-native:** Wrap v helper

```bash
# Vytvorit:
cmd/lint-helper/main.go

# Potom Makefile:
lint:
    @echo "Linting..."
    go run ./cmd/lint-helper
```

---

### Príkaz 3: Security Scan (Riadok 50)
**Aktuálne:** `if command -v govulncheck`
**Go-native:** Helper utility

```bash
# Vytvorit:
cmd/security-scanner/main.go

# Potom Makefile:
security-scan:
    @echo "Running security scans..."
    go run ./cmd/security-scanner
```

---

### Príkaz 4: Release (Riadok 67)
**Aktuálne:** `ls -lh | grep | awk`
**Go-native:** Existuje `cmd/release-manager/main.go` ✅

```bash
# Prepísať Makefile:
release: clean build
    @echo "Creating release artifacts..."
    go run ./cmd/release-manager -version $(VERSION) -bindir $(BIN_DIR) -outdir $(RELEASE_DIR)
```

---

### Príkaz 5: Backup Final (Riadok 71-74)
**Aktuálne:** `git archive | gzip | ls | awk`
**Go-native:** Existuje `cmd/backup-manager/main.go` ✅

```bash
# Prepísať Makefile:
backup-final:
    @echo "Creating final deterministic backup (gold master)"
    go run ./cmd/backup-manager -version $(VERSION) -outdir $(RELEASE_DIR)
```

---

## 🔗 LINKEDIN UMIESTNENIE

### Status: README.md JUŽ HAS LINKEDIN ✅

**Aktuálne (README.md, riadok 91):**
```markdown
LinkedIn: linkedin.com/in/kristian-kasnik-03056a377
```

### Odporúčané pridania:

#### 1. CONTRIBUTING.md (Nový oddiel na konci)
```markdown
## Authors & Contact

**Lead Maintainer:** Kristian Kasnik  
**LinkedIn:** [Kristian Kasnik](https://www.linkedin.com/in/kristian-kasnik-03056a377/)  
**Email:** itssafer@itssafer.org

For security issues, see [SECURITY.md](../SECURITY.md)
```

#### 2. docs/INDEX.md (Nová sekcia)
```markdown
## Community & Support

- **GitHub Issues:** [Report bugs or request features](https://github.com/ITSsafer-DevOps/N-Audit-Sentinel/issues)
- **LinkedIn:** [Connect with maintainers](https://www.linkedin.com/in/kristian-kasnik-03056a377/)
- **Security Contact:** See [SECURITY.md](../SECURITY.md)
```

#### 3. SECURITY.md (Nový oddiel na konci)
```markdown
## Security Contact & Responsible Disclosure

For security vulnerabilities, please contact:
- **LinkedIn:** [Kristian Kasnik](https://www.linkedin.com/in/kristian-kasnik-03056a377/)
- **Email:** security@itssafer.org (if available)

Do not file public issues for security vulnerabilities.
```

---

## 📊 CHÝBAJÚCE DIAGRAMY (Podľa Tech-Grade Priority)

### VYSOKÁ PRIORITA (Enterprise-grade)

#### 1. K8s Workload Deployment Topology
**Umiestnenie:** `DEPLOYMENT.md` (po "Architecture" sekcii)
**Format:** Mermaid graph
**Obsah:** Pod → Service → RBAC → KMS → Storage
**Priorita:** 🔴 HIGH (technical accuracy for DevOps)

#### 2. Cilium Policy Flow (Sequence Diagram)
**Umiestnenie:** `docs/ARCHITECTURE_DIAGRAMS.md`
**Format:** Mermaid sequenceDiagram
**Obsah:** TUI → PolicyGen → Cilium → Kernel L3/L7 enforcement
**Priorita:** 🔴 HIGH (security-critical)

#### 3. TUI State Machine
**Umiestnenie:** `MANUAL-TEST-GUIDE.md`
**Format:** Mermaid stateDiagram-v2
**Obsah:** Banner → ScopeCollection → ScopeValidation → PolicyGen → ShellStart → CommandCapture → SealGeneration → Cleanup
**Priorita:** 🔴 HIGH (user workflow documentation)

### STREDNÁ PRIORITA

#### 4. Integration Test Flow
**Umiestnenie:** `TESTING_AND_VERIFICATION.md`
**Format:** Mermaid flowchart
**Obsah:** Setup → Deploy → Connect → Scope → Exec → Capture → Verify → Cleanup
**Priorita:** 🟠 MEDIUM (testing documentation)

#### 5. Error Handling & Recovery
**Umiestnenie:** `SECURITY.md` (nový oddiel)
**Format:** Mermaid graph
**Obsah:** Policy Error → Log → Cleanup, Scope Error → Retry, Seal Error → Fallback
**Priorita:** 🟠 MEDIUM (operational resilience)

---

## 📈 SUMMARY TABLE

| Kategória | Položky | Status |
|-----------|---------|--------|
| Súbory na zmazanie | 6 | 📋 Hotovo |
| Súbory na opravenie (lokalizácia) | 1 | ✏️ 2 zmeny |
| Makefile príkazy | 7 | ⚙️ Reformovateľné |
| LinkedIn umiestnenie | 3 | 🔗 3 miesta |
| Chýbajúce diagramy | 5 | 📊 Odporúčané |

---

## 🚀 ODPORÚČANÝ POSTUP IMPLEMENTÁCIE

### Fáza 1: Cleanup (2 dni)
```
1. Zmazať 6 redundantných súborov ❌
2. Opraviť 2 slovenčiny v COMPREHENSIVE_ENHANCEMENT_REPORT.md ✏️
3. Commit: "chore: cleanup duplicate and localized documentation"
```

### Fáza 2: Makefile Refactor (5-7 dni)
```
1. Vytvorit cmd/verify-deps/main.go
2. Vytvorit cmd/lint-helper/main.go
3. Vytvorit cmd/security-scanner/main.go
4. Prepísať Makefile targety (release, backup-final, verify-deps, lint, security-scan)
5. Test: make help, make verify-deps, make release, make backup-final
6. Commit: "refactor(makefile): convert bash to Go-native commands"
```

### Fáza 3: LinkedIn & Contact (1 deň)
```
1. Pridať LinkedIn section do CONTRIBUTING.md
2. Pridať LinkedIn section do docs/INDEX.md
3. Pridať security contact do SECURITY.md
4. Commit: "docs: add LinkedIn and contact information"
```

### Fáza 4: Diagram Enhancements (3-5 dni)
```
1. Vytvorit K8s Workload Deployment diagram (DEPLOYMENT.md)
2. Vytvorit Cilium Policy Flow diagram (docs/ARCHITECTURE_DIAGRAMS.md)
3. Vytvorit TUI State Machine (MANUAL-TEST-GUIDE.md)
4. Vytvorit Integration Test Flow (TESTING_AND_VERIFICATION.md)
5. Vytvorit Error Handling diagram (SECURITY.md)
6. Validate: `grep -r "stateDiagram\|sequenceDiagram" *.md`
7. Commit: "docs: add enterprise-grade Mermaid diagrams"
```

---

**Dátum:** 2025-12-11  
**Čas:** komprehenzívny audit kompletný
