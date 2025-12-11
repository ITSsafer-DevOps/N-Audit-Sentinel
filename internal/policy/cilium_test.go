package policy

import (
	"testing"
)

func TestGeneratePolicyObject_Basic(t *testing.T) {
	c := &CiliumClient{}
	podLabels := map[string]string{"app": "n-audit-sentinel"}
	infraDNS := []string{"1.2.3.4"}
	infraAPI := "10.0.0.5:6443"
	targetIPs := []string{"192.0.2.1"}
	targetDomains := []string{"example.com"}

	pol := c.generatePolicyObject("test-policy", "default", podLabels, infraDNS, infraAPI, targetIPs, targetDomains)
	if pol.ObjectMeta.Name != "test-policy" {
		t.Fatalf("unexpected name: %s", pol.ObjectMeta.Name)
	}
	if pol.ObjectMeta.Namespace != "default" {
		t.Fatalf("unexpected namespace: %s", pol.ObjectMeta.Namespace)
	}
	if pol.Spec == nil {
		t.Fatalf("spec is nil")
	}
	if len(pol.Spec.Egress) == 0 {
		t.Fatalf("expected egress rules, got none")
	}
}
