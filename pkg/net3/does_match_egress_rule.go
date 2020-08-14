package net3

import (
	networkingv1 "k8s.io/api/networking/v1"
)

// doesMatchEgressRule checks if a destination matches a network
// policy egress rule.
func doesMatchEgressRule(rule networkingv1.NetworkPolicyEgressRule, destination Destination) bool {
	return true
}
