#!/usr/bin/env bash
set -euo pipefail

echo "=========================================="
echo "N-Audit Sentinel Local Deployment Test"
echo "=========================================="
echo ""

# Phase 1: Prerequisites Check
echo "[1/6] Checking prerequisites..."
command -v docker >/dev/null 2>&1 || { echo "ERROR: docker not found in PATH"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "ERROR: kubectl not found in PATH"; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "ERROR: terraform not found in PATH"; exit 1; }
echo "✓ All prerequisites available"
echo ""

# Phase 2: Build Docker Image
echo "[2/6] Building Docker image..."
docker build -t n-audit-sentinel:local-test . || { echo "ERROR: Docker build failed"; exit 1; }
echo "✓ Docker image built successfully"
echo ""

# Phase 3: Load Image into Local Cluster
echo "[3/6] Loading image into local cluster..."
echo "Which local Kubernetes cluster are you using?"
echo "  1) minikube"
echo "  2) k3d"
echo "  3) kind"
echo "  4) skip (image already loaded)"
read -p "Enter choice (1-4): " cluster_choice

case $cluster_choice in
  1)
    echo "Loading into minikube..."
    minikube image load n-audit-sentinel:local-test || { echo "ERROR: Failed to load image"; exit 1; }
    ;;
  2)
    echo "Loading into k3d..."
    k3d image import n-audit-sentinel:local-test || { echo "ERROR: Failed to load image"; exit 1; }
    ;;
  3)
    echo "Loading into kind..."
    kind load docker-image n-audit-sentinel:local-test || { echo "ERROR: Failed to load image"; exit 1; }
    ;;
  4)
    echo "Skipping image load..."
    ;;
  *)
    echo "ERROR: Invalid choice"
    exit 1
    ;;
esac
echo "✓ Image loaded"
echo ""

# Phase 4: Deploy with Terraform
echo "[4/6] Deploying with Terraform..."
cd deploy/terraform

# Create terraform.tfvars
cat > terraform.tfvars <<EOF
namespace       = "default"
image_name      = "n-audit-sentinel"
image_tag       = "local-test"
pvc_storage_size = "1Gi"
EOF
echo "✓ terraform.tfvars created"

# Initialize and apply
terraform init || { echo "ERROR: terraform init failed"; exit 1; }
terraform apply -auto-approve || { echo "ERROR: terraform apply failed"; exit 1; }
echo "✓ Terraform deployment successful"
cd ../..
echo ""

# Phase 5: Wait for Pod and Interactive Test
echo "[5/6] Waiting for pod to be ready..."
kubectl wait --for=condition=Ready pod/n-audit-sentinel --timeout=60s || { echo "ERROR: Pod failed to become ready"; exit 1; }
echo "✓ Pod is ready"
echo ""

echo "=========================================="
echo "INTERACTIVE TEST"
echo "=========================================="
echo "The script will now attach to the pod."
echo "Please test the following:"
echo "  - Verify the banner appears"
echo "  - Enter pentester name and client name"
echo "  - Enter scope (IPs and domains) using double-enter"
echo "  - Test the safety loop by typing 'exit' (shell should respawn)"
echo ""
echo "To detach from the pod, press: Ctrl+P then Ctrl+Q"
echo ""
read -p "Press ENTER to attach to the pod..." 

kubectl attach -it n-audit-sentinel || echo "Detached from pod"
echo ""

echo "=========================================="
echo "GRACEFUL SHUTDOWN TEST"
echo "=========================================="
echo "The script will now trigger a graceful shutdown using the CLI tool."
read -p "Press ENTER to trigger shutdown..." 

kubectl exec -it n-audit-sentinel -- /usr/local/bin/n-audit || echo "Pod terminated"
echo "✓ Graceful shutdown completed"
echo ""

# Phase 6: Cleanup
echo "[6/6] Cleaning up..."
cd deploy/terraform
terraform destroy -auto-approve || { echo "WARNING: terraform destroy had issues"; }
cd ../..
echo "✓ Cleanup complete"
echo ""

echo "=========================================="
echo "Local deployment test complete."
echo "=========================================="
