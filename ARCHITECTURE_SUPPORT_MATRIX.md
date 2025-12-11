# Architecture & Platform Support Matrix

## Supported Architectures

| Architecture | Status | Platform Support |
|--------------|--------|------------------|
| linux/amd64 | ✅ Tested | K3s, K8s, Talos, MicroShift, OpenShift |
| linux/arm64 | ✅ Ready | All (Go code is platform-agnostic) |
| linux/armv7 | ✅ Ready | All (Go code is platform-agnostic) |

## Kubernetes Distribution Support

| Distribution | Version | Status | Deployment | Tests |
|--------------|---------|--------|------------|-------|
| K3s | 1.19+ | ✅ Full | deploy/k3s/ | tests/e2e/k3s/ |
| Kubernetes | 1.19+ | ✅ Full | deploy/k8s/ | tests/e2e/k8s/ |
| Talos Linux | Latest | ✅ Full | deploy/talos/ | tests/e2e/talos/ |
| MicroShift | 4.x+ | ✅ Full | deploy/microshift/ | tests/e2e/microshift/ |
| OpenShift | 4.x+ | ✅ Full | deploy/openshift/ | tests/e2e/openshift/ |

## CNI Support

| CNI | Version | Status |
|-----|---------|--------|
| Cilium | 1.10+ | ✅ Required for policies |
| Calico | 3.x+ | ✅ Compatible |
| Weave | Latest | ✅ Compatible |

## Storage Support

| Storage | Type | Status |
|---------|------|--------|
| hostPath | Local | ✅ Supported |
| PersistentVolume | Cluster | ✅ Supported |
| NFS | Network | ✅ Supported |
| EBS | AWS | ✅ Supported |

## RBAC & Access Control

- ✅ ServiceAccount required
- ✅ ClusterRole for read access
- ✅ ClusterRoleBinding for binding
- ✅ Namespace scoping supported
