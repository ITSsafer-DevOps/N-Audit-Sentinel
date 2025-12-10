# N-Audit Sentinel Terraform Deployment

This directory contains production-ready Terraform configuration to deploy N-Audit Sentinel to Kubernetes.

Files:
- `variables.tf`: Input variables with sensible defaults.
- `main.tf`: PVC and Pod resources using the provided variables.
- `outputs.tf`: Convenience outputs (pod name).
- `terraform.tfvars.example`: Example variable values to copy and tailor.

For full, step-by-step instructions including image build/push and operational guidance, see the top-level deployment guide:

- `../../DEPLOYMENT.md`

Quick start:
1. Copy `terraform.tfvars.example` to `terraform.tfvars` and edit values.
2. Run `terraform init && terraform apply`.
3. Attach to the pod via `kubectl attach -it n-audit-sentinel -n <namespace>`.
