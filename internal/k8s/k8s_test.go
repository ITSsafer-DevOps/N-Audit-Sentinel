package k8s

import "testing"

func TestInitClient(t *testing.T) {
	c, err := InitClient("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil || c.Config != "foo" {
		t.Fatalf("unexpected client: %+v", c)
	}
}
