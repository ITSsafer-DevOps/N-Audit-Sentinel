package cilium

import (
	"strings"
	"testing"
)

func TestGenerateCiliumPolicy(t *testing.T) {
	name := "testpolicy"
	cidrs := []string{"1.2.3.4/32", "10.0.0.0/8"}
	out := GenerateCiliumPolicy(name, cidrs)
	if !strings.Contains(out, "name: "+name) {
		t.Fatalf("policy missing name: %s", out)
	}
	for _, c := range cidrs {
		if !strings.Contains(out, "\""+c+"\"") {
			t.Fatalf("policy missing cidr %s: %s", c, out)
		}
	}
}
