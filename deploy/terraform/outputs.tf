output "pod_name" {
  description = "Deployed N-Audit Sentinel pod name"
  value       = kubernetes_pod.n_audit.metadata[0].name
}
// Placeholder Terraform outputs.tf for N-Audit Sentinel
