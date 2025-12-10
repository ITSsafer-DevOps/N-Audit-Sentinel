variable "namespace" {
  description = "Kubernetes namespace to deploy into"
  type        = string
  default     = "default"
}

variable "image_name" {
  description = "Container image name (e.g., your-registry/n-audit-sentinel)"
  type        = string
  default     = "your-registry/n-audit-sentinel"
}

variable "image_tag" {
  description = "Container image tag (e.g., v1.0.0)"
  type        = string
  default     = "latest"
}

variable "pvc_storage_size" {
  description = "PersistentVolumeClaim storage size for logs/artifacts"
  type        = string
  default     = "50Gi"
}

variable "pvc_storage_class" {
  description = "StorageClass name for PVC (optional)"
  type        = string
  default     = ""
}
// Placeholder Terraform variables.tf for N-Audit Sentinel
