package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverK8sAPI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("KUBERNETES_SERVICE_HOST", "10.43.0.1")
		t.Setenv("KUBERNETES_SERVICE_PORT", "443")
		endpoint, err := DiscoverK8sAPI()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "10.43.0.1:443"
		if endpoint != expected {
			t.Fatalf("expected %q got %q", expected, endpoint)
		}
	})

	t.Run("missing_env_vars", func(t *testing.T) {
		t.Setenv("KUBERNETES_SERVICE_HOST", "")
		t.Setenv("KUBERNETES_SERVICE_PORT", "")
		_, err := DiscoverK8sAPI()
		if err == nil {
			t.Fatal("expected error when env vars not set, got nil")
		}
	})
}

func TestDiscoverDNS(t *testing.T) {
	tmpDir := t.TempDir()
	resolvPath := filepath.Join(tmpDir, "resolv.conf")

	content := `# This is a comment
nameserver 8.8.8.8
search my.cluster.local
nameserver 1.1.1.1
`
	if err := os.WriteFile(resolvPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp resolv.conf: %v", err)
	}

	nameservers, err := DiscoverDNS(resolvPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"8.8.8.8", "1.1.1.1"}
	if len(nameservers) != len(expected) {
		t.Fatalf("expected %d nameservers, got %d", len(expected), len(nameservers))
	}

	for i, ns := range nameservers {
		if ns != expected[i] {
			t.Errorf("nameserver[%d]: expected %q got %q", i, expected[i], ns)
		}
	}
}

func TestDiscoverDNS_FileNotFound(t *testing.T) {
	_, err := DiscoverDNS("/nonexistent/resolv.conf")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}
