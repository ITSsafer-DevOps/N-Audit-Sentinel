package policy

import (
	"strings"
	"testing"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

func TestGeneratePolicyObject_Basic(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"app": "test"}
	infraDNS := []string{"8.8.8.8", "8.8.4.4"}
	infraAPI := "10.0.0.1:443"
	targetIPs := []string{"192.168.1.0/24", "10.0.0.0"}
	targetDomains := []string{"target.com", "test.org"}

	policy := client.generatePolicyObject(
		"test-policy", "default",
		podLabels, infraDNS, infraAPI, targetIPs, targetDomains,
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if policy.Name != "test-policy" {
		t.Fatalf("expected policy name 'test-policy', got %s", policy.Name)
	}
	if policy.Namespace != "default" {
		t.Fatalf("expected namespace 'default', got %s", policy.Namespace)
	}
}

func TestGeneratePolicyObject_NoDNS(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"app": "scanner"}
	infraDNS := []string{}
	infraAPI := "10.1.0.1:6443"
	targetIPs := []string{"172.16.0.0/16"}
	targetDomains := []string{}

	policy := client.generatePolicyObject(
		"scanner-policy", "scanning",
		podLabels, infraDNS, infraAPI, targetIPs, targetDomains,
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if len(policy.Spec.Egress) == 0 {
		t.Fatalf("expected non-empty egress rules")
	}
}

func TestGeneratePolicyObject_EmptyTargets(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"tier": "audit"}
	infraDNS := []string{"1.1.1.1"}
	infraAPI := "10.2.0.1:443"
	targetIPs := []string{}
	targetDomains := []string{}

	policy := client.generatePolicyObject(
		"audit-policy", "audit-ns",
		podLabels, infraDNS, infraAPI, targetIPs, targetDomains,
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if policy.Spec == nil {
		t.Fatalf("expected non-nil Spec")
	}
}

func TestGeneratePolicyObject_NoInfra(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"env": "prod"}
	infraDNS := []string{}
	infraAPI := ""
	targetIPs := []string{"203.0.113.0/24"}
	targetDomains := []string{"prod.example.com"}

	policy := client.generatePolicyObject(
		"prod-policy", "prod",
		podLabels, infraDNS, infraAPI, targetIPs, targetDomains,
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
}

func TestGeneratePolicyObject_TypeAssertion(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"test": "true"}
	policy := client.generatePolicyObject(
		"type-test", "default",
		podLabels, []string{"8.8.8.8"}, "10.0.0.1:443",
		[]string{"10.0.0.0/8"}, []string{"test.com"},
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if _, ok := interface{}(policy).(*ciliumv2.CiliumNetworkPolicy); !ok {
		t.Fatalf("expected *CiliumNetworkPolicy type")
	}
}

func TestGeneratePolicyObject_EndpointSelector(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"app": "endpoint-test", "tier": "web"}
	policy := client.generatePolicyObject(
		"endpoint-policy", "default",
		podLabels, []string{}, "", []string{}, []string{},
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if policy.Spec == nil || policy.Spec.EndpointSelector.LabelSelector == nil {
		t.Fatalf("expected EndpointSelector with LabelSelector")
	}
	if policy.Spec.EndpointSelector.LabelSelector.MatchLabels["app"] != "endpoint-test" {
		t.Fatalf("expected pod label 'app: endpoint-test'")
	}
	if policy.Spec.EndpointSelector.LabelSelector.MatchLabels["tier"] != "web" {
		t.Fatalf("expected pod label 'tier: web'")
	}
}

func TestGeneratePolicyObject_EgressRulesStructure(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{"app": "test"}
	infraDNS := []string{"8.8.8.8"}
	infraAPI := "10.0.0.1:443"
	targetIPs := []string{"192.168.0.0/16"}
	targetDomains := []string{"test.com"}

	policy := client.generatePolicyObject(
		"order-test", "default",
		podLabels, infraDNS, infraAPI, targetIPs, targetDomains,
	)

	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if policy.Spec == nil || len(policy.Spec.Egress) == 0 {
		t.Fatalf("expected egress rules")
	}
	foundMaintenanceRule := false
	for _, rule := range policy.Spec.Egress {
		if len(rule.ToFQDNs) > 0 {
			for _, fqdn := range rule.ToFQDNs {
				if strings.Contains(fqdn.MatchPattern, "kali") {
					foundMaintenanceRule = true
					break
				}
			}
		}
	}
	if !foundMaintenanceRule {
		t.Fatalf("expected maintenance rule with kali.org in FQDN patterns")
	}
}

func TestNewCiliumClient_Initialization(t *testing.T) {
	// Test that basic struct initialization works
	client := &CiliumClient{}
	if client == nil {
		t.Fatalf("expected non-nil CiliumClient")
	}
}

func TestGeneratePolicyObject_EmptyPodLabels(t *testing.T) {
	client := &CiliumClient{}
	podLabels := map[string]string{}
	policy := client.generatePolicyObject(
		"empty-labels", "ns",
		podLabels, []string{}, "", []string{}, []string{},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
}

func TestGeneratePolicyObject_LongStrings(t *testing.T) {
	client := &CiliumClient{}
	longDomain := "very.very.very.long.subdomain.chain.example.company.infrastructure.service.com"
	longIP := "192.0.2.0/24"
	policy := client.generatePolicyObject(
		"long-test", "long-ns",
		map[string]string{"x": "y"},
		[]string{"8.8.8.8", "8.8.4.4"},
		"10.0.0.1:443",
		[]string{longIP},
		[]string{longDomain},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
}

func TestGeneratePolicyObject_ManyDNSServers(t *testing.T) {
	client := &CiliumClient{}
	dnsList := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1", "208.67.222.222"}
	policy := client.generatePolicyObject(
		"many-dns", "ns",
		map[string]string{"role": "dns"},
		dnsList,
		"10.0.0.1:443",
		[]string{},
		[]string{},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if len(policy.Spec.Egress) == 0 {
		t.Fatalf("expected egress rules for DNS servers")
	}
}

func TestGeneratePolicyObject_AllComponentsPresent(t *testing.T) {
	client := &CiliumClient{}
	policy := client.generatePolicyObject(
		"full-spec", "production",
		map[string]string{"app": "scanner", "tier": "security"},
		[]string{"8.8.8.8"},
		"10.0.0.1:443",
		[]string{"10.0.0.0/8", "172.16.0.0/12"},
		[]string{"internal.example.com", "external.target.org"},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if policy.Spec == nil || policy.Spec.EndpointSelector.LabelSelector == nil {
		t.Fatalf("expected complete policy spec")
	}
	// Verify multiple label selectors
	labels := policy.Spec.EndpointSelector.LabelSelector.MatchLabels
	if labels["app"] != "scanner" || labels["tier"] != "security" {
		t.Fatalf("expected all pod labels in policy")
	}
}

func TestGeneratePolicyObject_InfraAPIValidPort(t *testing.T) {
	client := &CiliumClient{}
	policy := client.generatePolicyObject(
		"api-port-test", "default",
		map[string]string{"app": "api"},
		[]string{"8.8.8.8"},
		"10.0.0.1:6443",
		[]string{},
		[]string{},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	foundAPIRule := false
	for _, rule := range policy.Spec.Egress {
		if len(rule.ToPorts) > 0 {
			for _, pr := range rule.ToPorts {
				for _, p := range pr.Ports {
					if p.Port == "6443" {
						foundAPIRule = true
						break
					}
				}
			}
		}
	}
	if !foundAPIRule {
		t.Fatalf("expected API rule with port 6443")
	}
}

func TestGeneratePolicyObject_MixedCIDRNotation(t *testing.T) {
	client := &CiliumClient{}
	targetIPs := []string{
		"192.0.2.1",     // single IP (should become /32)
		"10.0.0.0/24",   // CIDR with /24
		"172.16.0.0/12", // CIDR with /12
	}
	policy := client.generatePolicyObject(
		"cidr-test", "default",
		map[string]string{"tier": "network"},
		[]string{},
		"",
		targetIPs,
		[]string{},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	if len(policy.Spec.Egress) == 0 {
		t.Fatalf("expected egress rules for target IPs")
	}
}

func TestGeneratePolicyObject_MaintenanceWhitelist(t *testing.T) {
	client := &CiliumClient{}
	policy := client.generatePolicyObject(
		"maintenance-check", "default",
		map[string]string{"type": "maintenance"},
		[]string{},
		"",
		[]string{},
		[]string{},
	)
	if policy == nil {
		t.Fatalf("expected non-nil policy")
	}
	foundMaintenance := false
	requiredDomains := map[string]bool{
		"kali.org":   false,
		"github.com": false,
		"docker.io":  false,
		"gitlab.com": false,
		"pypi.org":   false,
		"crates.io":  false,
	}
	for _, rule := range policy.Spec.Egress {
		if len(rule.ToFQDNs) > 0 {
			for _, fqdn := range rule.ToFQDNs {
				for domain := range requiredDomains {
					if strings.Contains(fqdn.MatchPattern, domain) {
						requiredDomains[domain] = true
						foundMaintenance = true
					}
				}
			}
		}
	}
	if !foundMaintenance {
		t.Fatalf("expected maintenance whitelist domains in policy")
	}
}
