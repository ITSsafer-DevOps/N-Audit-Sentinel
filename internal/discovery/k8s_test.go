package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverK8sAPI(t *testing.T) {
	oldHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	oldPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	defer os.Setenv("KUBERNETES_SERVICE_HOST", oldHost)
	defer os.Setenv("KUBERNETES_SERVICE_PORT", oldPort)

	os.Setenv("KUBERNETES_SERVICE_HOST", "10.0.0.1")
	os.Setenv("KUBERNETES_SERVICE_PORT", "6443")

	ep, err := DiscoverK8sAPI()
	if err != nil {
		t.Fatalf("DiscoverK8sAPI error: %v", err)
	}
	if ep != "10.0.0.1:6443" {
		t.Fatalf("unexpected endpoint: %s", ep)
	}
}

func TestDiscoverDNS(t *testing.T) {
	dir := t.TempDir()
	rc := filepath.Join(dir, "resolv.conf")
	content := `# sample
nameserver 1.1.1.1
nameserver 8.8.8.8

# comment
`
	if err := os.WriteFile(rc, []byte(content), 0o644); err != nil {
		t.Fatalf("write resolv: %v", err)
	}
	ns, err := DiscoverDNS(rc)
	if err != nil {
		t.Fatalf("DiscoverDNS error: %v", err)
	}
	if len(ns) != 2 || ns[0] != "1.1.1.1" || ns[1] != "8.8.8.8" {
		t.Fatalf("unexpected nameservers: %v", ns)
	}
}
