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
