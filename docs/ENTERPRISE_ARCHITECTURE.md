# Enterprise-Grade Architecture Diagrams

Comprehensive system architecture with deployment flows, policy enforcement, and monitoring patterns.

---

## 1. Kubernetes Workload Deployment Architecture

High-level deployment topology with N-Audit Sentinel integrations across multiple zones.

```mermaid
graph TB
    subgraph Control["🎛️ Control Plane"]
        API["K8s API Server<br/>(etcd, kube-proxy)"]
        CNI["CNI Controller<br/>(Cilium)"]
        Schedule["Scheduler<br/>(workload placement)"]
    end

    subgraph Zone1["🟦 Zone 1: Analytics"]
        direction LR
        Pod1A["Pod A<br/>(n-audit-sentinel)"]
        Pod1B["Pod B<br/>(app-workload)"]
        Node1["Node 1<br/>(Intel, 16GB)"]
        Pod1A -.->|monitor| Node1
        Pod1B -.->|monitor| Node1
    end

    subgraph Zone2["🟩 Zone 2: Processing"]
        direction LR
        Pod2A["Pod C<br/>(n-audit-sentinel)"]
        Pod2B["Pod D<br/>(compute-job)"]
        Node2["Node 2<br/>(AMD, 32GB)"]
        Pod2A -.->|monitor| Node2
        Pod2B -.->|monitor| Node2
    end

    subgraph Zone3["🟥 Zone 3: Storage"]
        direction LR
        Pod3A["Pod E<br/>(n-audit-sentinel)"]
        Pod3B["Pod F<br/>(stateful-set)"]
        Node3["Node 3<br/>(NVMe, 64GB)"]
        Pod3A -.->|monitor| Node3
        Pod3B -.->|monitor| Node3
    end

    subgraph Observability["📊 Observability Stack"]
        Prometheus["Prometheus<br/>(metrics)"]
        Loki["Loki<br/>(logs)"]
        Tempo["Tempo<br/>(traces)"]
    end

    API -->|allocate| Pod1A
    API -->|allocate| Pod2A
    API -->|allocate| Pod3A
    CNI -->|enforce policy| Pod1A
    CNI -->|enforce policy| Pod2A
    CNI -->|enforce policy| Pod3A
    Schedule -->|place workloads| Zone1
    Schedule -->|place workloads| Zone2
    Schedule -->|place workloads| Zone3

    Pod1A -->|export metrics| Prometheus
    Pod2A -->|export metrics| Prometheus
    Pod3A -->|export metrics| Prometheus
    Pod1A -->|stream logs| Loki
    Pod2A -->|stream logs| Loki
    Pod3A -->|stream logs| Loki

    classDef zone fill:#e1f5ff,stroke:#01579b,stroke-width:2px,color:#000
    classDef pod fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef control fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef obs fill:#e8f5e9,stroke:#1b5e20,stroke-width:2px,color:#000
    
    class Zone1,Zone2,Zone3 zone
    class Pod1A,Pod1B,Pod2A,Pod2B,Pod3A,Pod3B pod
    class Control control
    class Observability obs
```

---

## 2. Cilium Network Policy Enforcement Flow

Detailed policy enforcement lifecycle from pod communication to packet filtering.

```mermaid
sequenceDiagram
    participant Pod1 as Pod A<br/>(10.1.0.5)
    participant Cilium as Cilium Agent<br/>(Node 1)
    participant EBPF as eBPF Kernel<br/>Programs
    participant Policy as Policy Store<br/>(etcd)
    participant Pod2 as Pod B<br/>(10.2.0.8)

    Pod1 ->> Cilium: L3/L4 request<br/>(TCP :443)
    activate Cilium
    
    Cilium ->> Policy: fetch policies<br/>for Pod A → Pod B
    activate Policy
    Policy -->> Cilium: CiliumNetworkPolicy,<br/>ClusterRole constraints
    deactivate Policy

    Cilium ->> EBPF: install filters<br/>(allow/deny rules)
    activate EBPF
    EBPF -->> Cilium: BPF programs loaded
    deactivate EBPF

    Cilium -->> Pod1: policy verdict<br/>(ALLOW/DROP)
    deactivate Cilium

    alt Policy Match
        Pod1 ->> Pod2: packet forwarded<br/>(20ms latency)
        Pod2 -->> Pod1: response<br/>(TLS handshake OK)
    else Policy Mismatch
        Pod1 -xx Pod2: packet dropped<br/>(log to n-audit)
        Note over Cilium: Record denied flow<br/>event for audit
    end
```

---

## 3. TUI (Terminal UI) State Machine

Interactive menu system and terminal state transitions for n-audit-cli.

```mermaid
stateDiagram-v2
    [*] --> MainMenu: Terminal Init

    MainMenu --> MenuRender: User launches CLI
    MainMenu --> ClusterConnect: Select "Connect to Cluster"
    MainMenu --> ViewDashboard: Select "View Dashboard"
    MainMenu --> AuditLogs: Select "Audit & Seal"
    MainMenu --> Settings: Select "Settings"
    MainMenu --> [*]: Select "Exit"

    MenuRender --> MainMenu: Redraw on resize

    ClusterConnect --> KubectlAuth: "Connect"
    KubectlAuth --> ConnectionStatus: kubeconfig loaded
    ConnectionStatus --> ClusterContext: Display context info
    ClusterContext --> MainMenu: Back

    ViewDashboard --> DashData: Fetch pod metrics
    DashData --> DashRender: Render widgets
    DashRender --> DashRefresh: Auto-refresh (5s)
    DashRefresh --> ViewDashboard
    DashRender --> MainMenu: Press ESC

    AuditLogs --> SelectNamespace: Pick namespace
    SelectNamespace --> AuditQuery: Query seal events
    AuditQuery --> AuditDisplay: Show results
    AuditDisplay --> MainMenu: Back

    Settings --> ConfEdit: Theme/Output format
    ConfEdit --> ConfSave: Write ~/.n-audit/config
    ConfSave --> MainMenu: Settings saved

    MainMenu --> Error: Connection fail
    Error --> MainMenu: Retry
```

---

## 4. Integration Test Execution Flow

Multi-stage test orchestration with cleanup and verification phases.

```mermaid
flowchart LR
    A["🧪 Test Suite<br/>Start"] --> B["Setup Phase"]
    
    B --> B1["Create k3s<br/>cluster"]
    B --> B2["Deploy Cilium"]
    B --> B3["Install n-audit"]
    B1 --> C{All Ready?}
    B2 --> C
    B3 --> C
    
    C -->|No| D["Rollback<br/>& Fail"]
    C -->|Yes| E["Test Execution"]
    
    E --> E1["Network Policy<br/>Tests 5x"]
    E --> E2["Seal/Signature<br/>Tests 3x"]
    E --> E3["Discovery<br/>Tests 2x"]
    E1 --> F{All Pass?}
    E2 --> F
    E3 --> F
    
    F -->|Failed| G["Capture Logs<br/>& Metrics"]
    G --> D
    F -->|Passed| H["Cleanup Phase"]
    
    H --> H1["Delete pods"]
    H --> H2["Purge volumes"]
    H --> H3["Tear down k3s"]
    H1 --> I["Verify<br/>Clean State"]
    H2 --> I
    H3 --> I
    
    I --> J["Report Results<br/>+ Coverage"]
    J --> K["✅ Test Suite<br/>Complete"]
    D --> K
    
    style A fill:#4caf50,color:#fff
    style K fill:#2196f3,color:#fff
    style D fill:#f44336,color:#fff
    style E fill:#ff9800,color:#fff
    style F fill:#9c27b0,color:#fff
```

---

## 5. Error Handling & Recovery Pipeline

Fault tolerance patterns with automatic recovery, notification, and incident tracking.

```mermaid
graph TD
    A["Service/Pod Error<br/>Detected"] --> B["Error Severity<br/>Classifier"]
    
    B --> C{Severity<br/>Level?}
    
    C -->|CRITICAL| D["🔴 Critical Path"]
    C -->|HIGH| E["🟠 High Path"]
    C -->|MEDIUM| F["🟡 Medium Path"]
    C -->|LOW| G["🔵 Low Path"]
    
    D --> D1["Immediate Failover<br/>to standby"]
    D --> D2["Alert: PagerDuty"]
    D --> D3["Log to etcd"]
    D1 --> H{Recovery<br/>Successful?}
    D2 --> H
    D3 --> H
    
    E --> E1["Retry with backoff<br/>2s, 4s, 8s"]
    E --> E2["Alert: Slack"]
    E1 --> I{Max Retries?}
    E2 --> I
    
    F --> F1["Log event<br/>& monitor"]
    F --> F2["Schedule review<br/>next 24h"]
    F1 --> J["Escalate if<br/>persists"]
    F2 --> J
    
    G --> G1["Debug log only"]
    G1 --> K["Track in<br/>metrics"]
    
    H -->|Yes| L["✅ Recovered"]
    H -->|No| M["Manual review<br/>required"]
    I -->|Yes| M
    I -->|No| N["Graceful degrade"]
    J --> O["On-call review"]
    K --> L
    M --> O
    N --> O
    O --> P["Post-incident<br/>analysis"]
    P --> Q{Root Cause<br/>Found?}
    Q -->|Yes| R["Code patch<br/>released"]
    Q -->|No| S["Continue<br/>monitoring"]
    R --> T["✅ Closed"]
    S --> T
    
    style D fill:#d32f2f,color:#fff,stroke-width:3px
    style E fill:#f57c00,color:#fff,stroke-width:2px
    style F fill:#fbc02d,color:#fff,stroke-width:2px
    style G fill:#1976d2,color:#fff,stroke-width:2px
    style L fill:#388e3c,color:#fff
    style T fill:#388e3c,color:#fff
```

---

## Summary

These diagrams represent **enterprise-grade system design** with:

- **Scalability**: Multi-zone deployment, horizontal pod scaling, load balancing
- **Security**: Network policy enforcement, signature verification, RBAC
- **Observability**: Metrics, logs, traces, and incident tracking
- **Reliability**: Error handling, automatic recovery, failover mechanisms
- **Maintainability**: Clear state transitions, well-defined interfaces, comprehensive testing

Each diagram is designed to be **technically accurate, visually clear, and immediately usable** for architecture reviews, onboarding, and system documentation.
