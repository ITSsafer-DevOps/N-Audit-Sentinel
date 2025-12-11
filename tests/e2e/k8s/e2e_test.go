package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"testing"
)

func TestK8sAPIAccess(t *testing.T) {
	// E2E test runs only when RUN_E2E env var is set
	if os.Getenv("RUN_E2E") != "true" {
		t.Skip("skipping e2e test; set RUN_E2E=true to run")
	}
	cfg, err := rest.InClusterConfig()
	if err != nil {
		// fallback to default config using kubeconfig from env
		cfg, err = rest.InClusterConfig()
		if err != nil {
			t.Fatalf("could not get cluster config: %v", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create clientset: %v", err)
	}
	_, err = clientset.ServerVersion()
	if err != nil {
		t.Fatalf("could not query server version: %v", err)
	}
}
