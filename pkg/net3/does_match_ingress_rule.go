package net3

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// doesMatchIngressRule checks if a destination matches a network
// policy ingress rule.
func (n *net3) doesMatchIngressRule(rule networkingv1.NetworkPolicyIngressRule, src *v1.Pod, port int32) (bool, error) {
	doesPortMatch := false
	for _, p := range rule.Ports {
		if p.Port.IntVal == port {
			doesPortMatch = true
			break
		}
	}
	if !doesPortMatch {
		return false, nil
	}

	for _, from := range rule.From {
		if from.PodSelector != nil {
			if !doesMatchSelector(from.PodSelector.MatchLabels, src.Labels) {
				continue
			}
		}
		if from.NamespaceSelector != nil {
			ns, err := n.k8s.CoreV1().Namespaces().Get(context.Background(), src.Namespace, metav1.GetOptions{})
			if err != nil {
				return false, fmt.Errorf("error getting source namespace %q: %w", src.Namespace, err)
			}
			if !doesMatchSelector(from.NamespaceSelector.MatchLabels, ns.Labels) {
				continue
			}
		}
		if from.IPBlock != nil {
			doesMatch, err := doesMatchIPBlock(*from.IPBlock, src.Status.PodIP)
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
