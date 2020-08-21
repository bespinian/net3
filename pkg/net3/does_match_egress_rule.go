package net3

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// doesMatchEgressRule checks if a destination matches a network
// policy egress rule.
func (n *net3) doesMatchEgressRule(rule networkingv1.NetworkPolicyEgressRule, dest *v1.Pod, port int32, namespace string) (bool, error) {
	doesPortMatch := false
	if len(rule.Ports) == 0 {
		doesPortMatch = true
	}
	for _, p := range rule.Ports {
		if p.Port.IntVal == port {
			doesPortMatch = true
			break
		}
	}
	if !doesPortMatch {
		return false, nil
	}

	for _, from := range rule.To {
		if from.PodSelector != nil {
			if !doesMatchSelector(from.PodSelector.MatchLabels, dest.Labels) {
				continue
			}
		}
		if from.NamespaceSelector != nil {
			ns, err := n.k8s.CoreV1().Namespaces().Get(context.Background(), dest.Namespace, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("error getting destination namespace %q: %w", dest.Namespace, err)
			}
			if !doesMatchSelector(from.NamespaceSelector.MatchLabels, ns.Labels) {
				continue
			}
		} else if dest.Namespace != namespace {
			continue
		}
		if from.IPBlock != nil {
			doesMatch, err := doesMatchIPBlock(*from.IPBlock, dest.Status.PodIP)
			if err != nil {
				return false, fmt.Errorf("error checking if IP block matches: %w", err)
			}
			if !doesMatch {
				continue
			}
		}
		return true, nil
	}

	return false, nil
}
