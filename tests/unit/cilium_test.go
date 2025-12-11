package unit

import (
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/cilium"
	"strings"
	"testing"
)

func TestGenerateCiliumPolicy(t *testing.T) {
	yaml := cilium.GenerateCiliumPolicy("test-policy", []string{"10.0.0.0/8", "192.168.0.0/16"})
	if !strings.Contains(yaml, "CiliumNetworkPolicy") {
		t.Fatalf("policy output missing kind")
	}
	if !strings.Contains(yaml, "10.0.0.0/8") {
		t.Fatalf("cidr missing")
	}
}
