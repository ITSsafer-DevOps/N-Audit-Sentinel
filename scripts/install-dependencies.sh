#!/usr/bin/env bash
set -euo pipefail

echo "=========================================="
echo "N-Audit Sentinel Dependencies Installer"
echo "=========================================="
echo "This will install: Docker, K3s, kubectl"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "ERROR: Please run as root (sudo bash install-dependencies.sh)"
    exit 1
fi

echo "[1/4] Installing Docker..."
if command -v docker >/dev/null 2>&1; then
    echo "✓ Docker already installed: $(docker --version)"
else
    # Install Docker
    apt-get update
    apt-get install -y ca-certificates curl gnupg lsb-release
    
    # Add Docker's official GPG key
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg
    
    # Set up the repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker Engine
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Add current user to docker group
    usermod -aG docker $SUDO_USER || true
    
    echo "✓ Docker installed successfully"
fi
echo ""

echo "[2/4] Installing K3s (lightweight Kubernetes)..."
if command -v k3s >/dev/null 2>&1; then
    echo "✓ K3s already installed"
else
    # Install K3s with Cilium CNI
    curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--flannel-backend=none --disable-network-policy --disable=traefik" sh -
    
    # Wait for K3s to be ready
    echo "Waiting for K3s to be ready..."
    sleep 10
    
    # Setup kubeconfig for non-root user
    mkdir -p /home/$SUDO_USER/.kube
    cp /etc/rancher/k3s/k3s.yaml /home/$SUDO_USER/.kube/config
    chown -R $SUDO_USER:$SUDO_USER /home/$SUDO_USER/.kube
    
    echo "✓ K3s installed successfully"
fi
echo ""

echo "[3/4] Installing Cilium CNI..."
# Install Cilium CLI
CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt)
CLI_ARCH=amd64
if [ "$(uname -m)" = "aarch64" ]; then CLI_ARCH=arm64; fi
curl -L --fail --remote-name-all https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-${CLI_ARCH}.tar.gz.sha256sum
tar xzvfC cilium-linux-${CLI_ARCH}.tar.gz /usr/local/bin
rm cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}

# Install Cilium
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
cilium install --wait

echo "✓ Cilium installed successfully"
echo ""

echo "[4/4] Verifying installation..."
echo ""
echo "Docker version:"
docker --version
echo ""
echo "kubectl version:"
kubectl version --client
echo ""
echo "K3s status:"
systemctl status k3s --no-pager | head -n 5
echo ""
echo "Cilium status:"
cilium status --wait
echo ""

echo "=========================================="
echo "Installation Complete!"
echo "=========================================="
echo ""
echo "⚠️  IMPORTANT: Log out and back in for Docker group changes to take effect"
echo ""
echo "Verify installation:"
echo "  docker ps"
echo "  kubectl get nodes"
echo "  kubectl get pods -A"
echo ""
echo "Next steps:"
echo "  cd /home/$SUDO_USER/VScode-server/N-Audit\ Sentinel"
echo "  ./local-deploy-test.sh"
echo ""
