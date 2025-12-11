package policy

import (
	"errors"
	"testing"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

func TestApplyPolicy_Success_WithHook(t *testing.T) {
	called := false
	c := &CiliumClient{
		createPolicyFunc: func(namespace string, policy *v2.CiliumNetworkPolicy) (*v2.CiliumNetworkPolicy, error) {
			called = true
			if namespace != "default" {
				t.Fatalf("unexpected namespace: %s", namespace)
			}
			if policy == nil || policy.ObjectMeta.Name != "hook-policy" {
				t.Fatalf("unexpected policy passed to hook")
			}
			return policy, nil
		},
	}

	err := c.ApplyPolicy("hook-policy", "default", map[string]string{"app": "x"}, nil, "", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("createPolicyFunc hook was not called")
	}
}

func TestApplyPolicy_Error_WithHook(t *testing.T) {
	c := &CiliumClient{
		createPolicyFunc: func(namespace string, policy *v2.CiliumNetworkPolicy) (*v2.CiliumNetworkPolicy, error) {
			return nil, errors.New("create-failed")
		},
	}

	err := c.ApplyPolicy("err-policy", "ns", nil, nil, "", nil, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDeletePolicy_Success_WithHook(t *testing.T) {
	called := false
	c := &CiliumClient{
		deletePolicyFunc: func(namespace, name string) error {
			called = true
			if namespace != "kube" || name != "to-delete" {
				t.Fatalf("unexpected args: %s %s", namespace, name)
			}
			return nil
		},
	}

	if err := c.DeletePolicy("to-delete", "kube"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("deletePolicyFunc hook was not called")
	}
}

func TestDeletePolicy_Error_WithHook(t *testing.T) {
	c := &CiliumClient{
		deletePolicyFunc: func(namespace, name string) error {
			return errors.New("delete-failed")
		},
	}

	if err := c.DeletePolicy("x", "y"); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
