# Architecture Diagrams

This document aggregates additional Mermaid diagrams for architecture, CI/CD, and internal package flows.

## CI/CD Pipeline

```mermaid
flowchart LR
  A["Commit to repo"] --> B["CI: lint & unit tests"]
  B --> C{coverage >= 75%}
  C -->|yes| D["Build release artifact"]
  C -->|no| E["Fail and notify PR"]
  D --> F["Build container image"]
  F --> G["Push image to registry"]
  G --> H["Deploy to staging"]
  H --> I["Run e2e tests"]
  I --> J["Promote to production"]

  style A fill:#4A90E2,color:#fff
  style B fill:#50E3C2,color:#000
  style D fill:#F5A623,color:#000
  style J fill:#7ED321,color:#000
```

## Package Data Flow (Discovery → Policy → Enforcement)

```mermaid
graph LR
  subgraph app["N-Audit Sentinel (pkg)"]
    Disc["discovery"] --> Pol["policy (releasemgr/internal/policy)"]
    Pol --> Sig["signature"]
    Disc --> Log["logger"]
    Rec["recorder"] --> Log
    Log --> Sig
  end

  style Disc fill:#50E3C2,color:#000
  style Pol fill:#F5A623,color:#000
  style Sig fill:#BD10E0,color:#fff
  style Log fill:#7ED321,color:#000
  style Rec fill:#9013FE,color:#fff
```

## Policy Generation Sequence

```mermaid
sequenceDiagram
  participant TUI
  participant Discovery
  participant PolicyEngine
  participant K8sAPI

  TUI->>Discovery: request cluster info
  Discovery-->>TUI: cluster & DNS
  TUI->>PolicyEngine: submit scope
  PolicyEngine->>K8sAPI: apply CiliumNetworkPolicy
  K8sAPI-->>PolicyEngine: 201 Created
  PolicyEngine-->>TUI: policy applied
```

## Forensic Seal Creation (Session End)

```mermaid
sequenceDiagram
  participant Recorder
  participant Hasher
  participant Signer
  participant Storage

  Recorder->>Hasher: stream session bytes
  Hasher-->>Signer: sha256 hash
  Signer->>Storage: append signature + seal
  Storage-->>Recorder: seal confirmed
```
