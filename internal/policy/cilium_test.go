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
func TestGeneratePolicyObject_MultipleTargets(t *testing.T) {
	c := &CiliumClient{}
	podLabels := map[string]string{"app": "test", "env": "prod"}
	infraDNS := []string{"1.2.3.4", "1.2.3.5"}
	infraAPI := "10.0.0.5:6443"
	targetIPs := []string{"192.0.2.1", "192.0.2.2", "10.0.0.0/8"}
	targetDomains := []string{"example.com", "test.org", "api.example.com"}

	pol := c.generatePolicyObject("multi-policy", "prod", podLabels, infraDNS, infraAPI, targetIPs, targetDomains)
	if pol.ObjectMeta.Name != "multi-policy" {
		t.Fatalf("unexpected name: %s", pol.ObjectMeta.Name)
	}
	if len(pol.Spec.Egress) < 2 {
		t.Fatalf("expected multiple egress rules, got %d", len(pol.Spec.Egress))
	}
}

func TestGeneratePolicyObject_EmptyTargets(t *testing.T) {
	c := &CiliumClient{}
	podLabels := map[string]string{"app": "test"}
	infraDNS := []string{}
	infraAPI := ""
	targetIPs := []string{}
	targetDomains := []string{}

	pol := c.generatePolicyObject("empty-policy", "default", podLabels, infraDNS, infraAPI, targetIPs, targetDomains)
	if pol.ObjectMeta.Name != "empty-policy" {
		t.Fatalf("unexpected name: %s", pol.ObjectMeta.Name)
	}
	if pol.Spec == nil {
		t.Fatalf("spec is nil")
	}
	// Should still create policy, even with empty targets
}

func TestGeneratePolicyObject_Labels(t *testing.T) {
	c := &CiliumClient{}
	podLabels := map[string]string{"app": "n-audit-sentinel", "component": "auditor"}
	infraDNS := []string{"8.8.8.8"}
	infraAPI := "api.kubernetes.svc:6443"
	targetIPs := []string{}
	targetDomains := []string{}

	pol := c.generatePolicyObject("labeled-policy", "kube-system", podLabels, infraDNS, infraAPI, targetIPs, targetDomains)
	if len(pol.Spec.EndpointSelector.MatchLabels) == 0 {
		t.Fatalf("expected pod labels to be set in selector")
	}
	if pol.Spec.EndpointSelector.MatchLabels["app"] != "n-audit-sentinel" {
		t.Fatalf("expected pod label 'app' to match 'n-audit-sentinel'")
	}
}
