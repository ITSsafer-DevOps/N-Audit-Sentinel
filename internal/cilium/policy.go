package cilium

import (
	"fmt"
	"strings"
)

// GenerateCiliumPolicy returns a minimal CiliumNetworkPolicy YAML as a string for the given name and CIDRs
func GenerateCiliumPolicy(name string, cidrs []string) string {
	cidrLines := []string{}
	for _, c := range cidrs {
		cidrLines = append(cidrLines, fmt.Sprintf("      - %q", c))
	}

	policy := fmt.Sprintf(`apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: %s
spec:
  endpointSelector: {}
  ingress:
  - fromCIDR:
%s
`, name, strings.Join(cidrLines, "\n"))
	return policy
}
