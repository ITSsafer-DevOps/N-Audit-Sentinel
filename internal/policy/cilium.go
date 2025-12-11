// N-Audit Sentinel - Cilium Policy Generator
// Developer: Kristian Kasnik
// Company: ITSsafer-DevOps
// License: MIT License
// Copyright (c) 2025 Kristian Kasnik, ITSsafer-DevOps
// Generates and manages CiliumNetworkPolicy with 3-zone enforcement.
// Maintenance Zone whitelist: *.kali.org, github.com, docker.io, gitlab.com, pypi.org, crates.io
package policy

import (
	"context"
	"fmt"
	"strings"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumclientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	slimmetav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	ciliumapi "github.com/cilium/cilium/pkg/policy/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CiliumClient manages Cilium network policy operations.
// It wraps the versioned Cilium clientset to create and delete CiliumNetworkPolicy resources.
type CiliumClient struct {
	clientset ciliumclientset.Interface
}

// NewCiliumClient creates a new Cilium client.
// It attempts to use in-cluster Kubernetes configuration and falls back to the default kubeconfig.
// Returns a CiliumClient or an error if client initialization fails.
func NewCiliumClient() (*CiliumClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build k8s config: %w", err)
		}
	}

	clientset, err := ciliumclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cilium clientset: %w", err)
	}

	return &CiliumClient{clientset: clientset}, nil
}

// generatePolicyObject constructs a CiliumNetworkPolicy with 3-zone enforcement.
// Zones:
// 1) Infrastructure: K8s API and DNS access.
// 2) Maintenance: limited outbound to whitelisted repositories via HTTP/HTTPS.
// 3) Target Scope: user-defined IPs and domains.
func (c *CiliumClient) generatePolicyObject(
	policyName, namespace string,
	podLabels map[string]string,
	infraDNS []string,
	infraAPI string,
	targetIPs []string,
	targetDomains []string,
) *ciliumv2.CiliumNetworkPolicy {
	policy := &ciliumv2.CiliumNetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: namespace,
		},
		Spec: &ciliumapi.Rule{
			EndpointSelector: ciliumapi.EndpointSelector{
				LabelSelector: &slimmetav1.LabelSelector{
					MatchLabels: podLabels,
				},
			},
		},
	}

	var egressRules []ciliumapi.EgressRule

	// Zone 1: Infrastructure (K8s API + DNS)
	if infraAPI != "" {
		parts := strings.Split(infraAPI, ":")
		if len(parts) == 2 {
			apiRule := ciliumapi.EgressRule{
				EgressCommonRule: ciliumapi.EgressCommonRule{
					ToCIDRSet: ciliumapi.CIDRRuleSlice{
						{Cidr: ciliumapi.CIDR(parts[0] + "/32")},
					},
				},
				ToPorts: ciliumapi.PortRules{
					{
						Ports: []ciliumapi.PortProtocol{
							{Port: parts[1], Protocol: ciliumapi.ProtoTCP},
						},
					},
				},
			}
			egressRules = append(egressRules, apiRule)
		}
	}

	if len(infraDNS) > 0 {
		var dnsCIDRSet ciliumapi.CIDRRuleSlice
		for _, dns := range infraDNS {
			dnsCIDRSet = append(dnsCIDRSet, ciliumapi.CIDRRule{
				Cidr: ciliumapi.CIDR(dns + "/32"),
			})
		}
		dnsRule := ciliumapi.EgressRule{
			EgressCommonRule: ciliumapi.EgressCommonRule{
				ToCIDRSet: dnsCIDRSet,
			},
			ToPorts: ciliumapi.PortRules{
				{
					Ports: []ciliumapi.PortProtocol{
						{Port: "53", Protocol: ciliumapi.ProtoUDP},
					},
				},
			},
		}
		egressRules = append(egressRules, dnsRule)
	}

	// Zone 2: Maintenance (whitelisted repos with HTTP/HTTPS only)
	maintenanceDomains := []string{
		"*.kali.org",
		"github.com",
		"docker.io",
		"gitlab.com",
		"pypi.org",
		"crates.io",
	}
	var fqdnRules ciliumapi.FQDNSelectorSlice
	for _, domain := range maintenanceDomains {
		fqdnRules = append(fqdnRules, ciliumapi.FQDNSelector{
			MatchPattern: domain,
		})
	}
	maintenanceRule := ciliumapi.EgressRule{
		ToFQDNs: fqdnRules,
		ToPorts: ciliumapi.PortRules{
			{
				Ports: []ciliumapi.PortProtocol{
					{Port: "80", Protocol: ciliumapi.ProtoTCP},
					{Port: "443", Protocol: ciliumapi.ProtoTCP},
				},
			},
		},
	}
	egressRules = append(egressRules, maintenanceRule)

	// Zone 3: Target Scope (user-defined IPs and domains - allow all)
	if len(targetIPs) > 0 {
		var targetCIDRSet ciliumapi.CIDRRuleSlice
		for _, ip := range targetIPs {
			if !strings.Contains(ip, "/") {
				ip = ip + "/32"
			}
			targetCIDRSet = append(targetCIDRSet, ciliumapi.CIDRRule{
				Cidr: ciliumapi.CIDR(ip),
			})
		}
		targetIPRule := ciliumapi.EgressRule{
			EgressCommonRule: ciliumapi.EgressCommonRule{
				ToCIDRSet: targetCIDRSet,
			},
		}
		egressRules = append(egressRules, targetIPRule)
	}

	if len(targetDomains) > 0 {
		var targetFQDNs ciliumapi.FQDNSelectorSlice
		for _, domain := range targetDomains {
			targetFQDNs = append(targetFQDNs, ciliumapi.FQDNSelector{
				MatchName: domain,
			})
		}
		targetDomainRule := ciliumapi.EgressRule{
			ToFQDNs: targetFQDNs,
		}
		egressRules = append(egressRules, targetDomainRule)
	}

	policy.Spec.Egress = egressRules
	return policy
}

// ApplyPolicy creates the CiliumNetworkPolicy in the cluster.
// Parameters include policy name, namespace, selector labels, discovered infra, and validated scope.
func (c *CiliumClient) ApplyPolicy(
	policyName, namespace string,
	podLabels map[string]string,
	infraDNS []string,
	infraAPI string,
	targetIPs []string,
	targetDomains []string,
) error {
	policy := c.generatePolicyObject(policyName, namespace, podLabels, infraDNS, infraAPI, targetIPs, targetDomains)
	_, err := c.clientset.CiliumV2().CiliumNetworkPolicies(namespace).Create(
		context.TODO(),
		policy,
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to apply policy: %w", err)
	}
	return nil
}

// DeletePolicy removes the CiliumNetworkPolicy from the cluster.
// Returns an error on failure.
func (c *CiliumClient) DeletePolicy(policyName, namespace string) error {
	err := c.clientset.CiliumV2().CiliumNetworkPolicies(namespace).Delete(
		context.TODO(),
		policyName,
		metav1.DeleteOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}
	return nil
}
