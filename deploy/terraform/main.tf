terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.27.0"
    }
  }
}

provider "kubernetes" {}

resource "kubernetes_persistent_volume_claim" "n_audit_pvc" {
  metadata {
    name      = "n-audit-pvc"
    namespace = var.namespace
  }
  spec {
    access_modes       = ["ReadWriteOnce"]
    storage_class_name = var.pvc_storage_class != "" ? var.pvc_storage_class : null
    resources {
      requests = {
        storage = var.pvc_storage_size
      }
    }
  }
}

resource "kubernetes_pod" "n_audit" {
  metadata {
    name      = "n-audit-sentinel"
    namespace = var.namespace
    labels = {
      app = "n-audit-sentinel"
    }
  }

  spec {
    security_context {
      run_as_user = 0
    }

    init_container {
      name              = "setup-storage"
      image             = "busybox:latest"
      image_pull_policy = "IfNotPresent"
      command           = ["sh", "-c", "mkdir -p /mnt/n-audit-data && chmod 777 /mnt/n-audit-data"]
      security_context {
        privileged = true
      }
      volume_mount {
        name       = "host-root"
        mount_path = "/mnt"
      }
    }

    container {
      name              = "sentinel"
      image             = "${var.image_name}:${var.image_tag}"
      image_pull_policy = "IfNotPresent"

      stdin = true
      tty   = true

      command = ["/usr/local/bin/n-audit-sentinel"]

      env {
        name = "POD_NAMESPACE"
        value_from {
          field_ref {
            field_path = "metadata.namespace"
          }
        }
      }

      env {
        name  = "SSH_SIGN_KEY_PATH"
        value = "/var/lib/n-audit/signing/id_ed25519"
      }

      volume_mount {
        name       = "data"
        mount_path = "/var/lib/n-audit"
      }

      resources {}
    }

    volume {
      name = "data"
      persistent_volume_claim {
        claim_name = kubernetes_persistent_volume_claim.n_audit_pvc.metadata[0].name
      }
    }

    volume {
      name = "host-root"
      host_path {
        path = "/"
        type = "Directory"
      }
    }
  }
}
// Placeholder Terraform main.tf for N-Audit Sentinel deployment
// NOTE: Actual implementation will define Pod, PVC without resource limits.
