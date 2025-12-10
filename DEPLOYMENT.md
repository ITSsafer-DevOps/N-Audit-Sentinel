# Deployment Guide (K3s/K8s)

This guide describes a minimal, reliable deployment flow using a hostPath for persistence and a ServiceAccount with RBAC for Cilium.

```mermaid
flowchart TB
	A[Build release (make)] --> B[Build runtime image]
	B --> C[Import/pull image to cluster]
	C --> D[Prepare hostPath + SSH key]
	D --> E[Apply SA + RBAC]
	E --> F[Apply Pod manifest]
	F --> G[Attach to TUI]
	G --> H[n-audit exit -> seal]
	H --> I[Artifacts on hostPath]
```

## Prerequisites
- Kubernetes cluster (K3s or K8s) with Cilium CNI
- `kubectl` and container build tool (Docker/nerdctl)
- Node host path for logs: `/mnt/n-audit-data`
- Optional: Terraform (for automated, reproducible deployments)

## 1) Build or Import the Image
On the node (K3s example):
```bash
# From repo root
make release VERSION=v1.0.0-final
cd beta-test-deployment
sudo docker build -t n-audit-sentinel:v1.0.0-final .
# Import to k3s containerd
sudo docker save n-audit-sentinel:v1.0.0-final | sudo k3s ctr images import -
```

If you use a registry, push and pull instead of importing.

### Terraform-driven deployment (optional)
From `deploy/terraform`:
```bash
terraform init
terraform apply -auto-approve \
	-var="namespace=default" \
	-var="image_name=n-audit-sentinel" \
	-var="image_tag=v1.0.0-final"
```
Terraform will create the pod with `app: n-audit-sentinel` label, required volumes, and service account bindings. You can still use hostPath or switch to a storage class via variables.

## 2) Prepare hostPath and SSH Key
```bash
sudo mkdir -p /mnt/n-audit-data/signing
sudo chmod 777 /mnt/n-audit-data
sudo ssh-keygen -t ed25519 -N "" -f /mnt/n-audit-data/signing/id_ed25519 -C "n-audit-sentinel@k8s"
sudo chmod 600 /mnt/n-audit-data/signing/id_ed25519
sudo chmod 644 /mnt/n-audit-data/signing/id_ed25519.pub
```

## 3) ServiceAccount + RBAC (Cilium CRDs)
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
	name: n-audit-sentinel
	namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
	name: n-audit-cilium-policy
rules:
- apiGroups: ["cilium.io"]
	resources: ["ciliumnetworkpolicies"]
	verbs: ["get", "list", "create", "delete", "update", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
	name: n-audit-cilium-policy
roleRef:
	apiGroup: rbac.authorization.k8s.io
	kind: ClusterRole
	name: n-audit-cilium-policy
subjects:
- kind: ServiceAccount
	name: n-audit-sentinel
	namespace: default
```
Apply it:
```bash
kubectl apply -f - <<'EOF'
# paste YAML above
EOF
```

Tip: This repo already includes ready-to-apply examples:
- `beta-test-deployment/serviceaccount.yaml`
- `beta-test-deployment/pod-fixed.yaml`

## 4) Pod Manifest (hostPath + labels)
```yaml
apiVersion: v1
kind: Pod
metadata:
	name: n-audit-sentinel
	labels:
		app: n-audit-sentinel
spec:
	serviceAccountName: n-audit-sentinel
	securityContext:
		runAsUser: 0
	containers:
	- name: sentinel
		image: n-audit-sentinel:v1.0.0-final
		stdin: true
		tty: true
		env:
		- name: SSH_SIGN_KEY_PATH
			value: "/var/lib/n-audit/signing/id_ed25519"
		volumeMounts:
		- name: data
			mountPath: /var/lib/n-audit
	volumes:
	- name: data
		hostPath:
			path: /mnt/n-audit-data
			type: DirectoryOrCreate
	restartPolicy: Always
```
Apply it:
```bash
kubectl apply -f pod.yaml
```

## 5) Attach and Operate
```bash
kubectl attach -it n-audit-sentinel -c sentinel
```
Follow TUI prompts (Pentester, Client, IP/CIDR scope, Domains). Use double‑Enter to finalize each list.

To end the session gracefully:
```bash
kubectl exec n-audit-sentinel -c sentinel -- /usr/local/bin/n-audit
```

## 6) Logs and Seal
- Pod log file: `/var/lib/n-audit/session.log`
- Host path: `/mnt/n-audit-data/session.log`
- Seal block with SHA‑256 + SSH signature is appended at teardown.

## Troubleshooting
- Policy apply timeout to `10.43.0.1:443`:
	- Ensure ServiceAccount + RBAC are applied and Cilium is healthy.
	- Verify the pod can reach the cluster IP of the API server.
- Missing seal: ensure private key exists and `SSH_SIGN_KEY_PATH` is set.
- Banner missed on attach: restart pod and attach immediately at start.

## See Also
- `VERIFICATION_GUIDE.md` for functional checks
- `README.md` for architecture and security model
