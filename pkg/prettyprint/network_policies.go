package prettyprint

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/networking/v1"
)

// NetworkPolicies prints multiple network policies.
func NetworkPolicies(policyType v1.PolicyType, policies []v1.NetworkPolicy, isAllowed bool) {
	for _, p := range policies {
		lines := []string{
			fmt.Sprintf("%s Network Policy", policyType),
			fmt.Sprintf("Name: %s", p.Name),
			fmt.Sprintf("Namespace: %s", p.Namespace),
		}

		var ruleStrings []string
		if policyType == v1.PolicyTypeIngress {
			for _, r := range p.Spec.Ingress {
				portStrings := make([]string, 0, len(r.Ports))
				for _, port := range r.Ports {
					portStrings = append(portStrings, fmtPort(port))
				}
				if len(portStrings) == 0 {
					portStrings = append(portStrings, "all traffic")
				}
				fromStrings := make([]string, 0, len(r.From))
				for _, from := range r.From {
					fromStrings = append(fromStrings, fmtPeer(from))
				}
				ruleStrings = append(ruleStrings, fmt.Sprintf("Allow %s from %s", strings.Join(portStrings, ", "), strings.Join(fromStrings, ", ")))
			}
		}
		for _, r := range ruleStrings {
			lines = append(lines, fmt.Sprintf("Rule: %s", r))
		}

		lines = append(lines, fmt.Sprintf("Allowing: %v", isAllowed))

		fmt.Print(asBox(lines))
	}
}

func fmtPort(port v1.NetworkPolicyPort) string {
	if port.Protocol == nil && port.Port == nil {
		return "all TCP traffic"
	}
	if port.Port == nil {
		return fmt.Sprintf("all %s traffic", *port.Protocol)
	}
	if port.Port.IntVal == 0 && port.Port.StrVal == "" {
		return fmt.Sprintf("all %s traffic", *port.Protocol)
	}
	if port.Port.IntVal == 0 {
		return port.Port.StrVal
	}
	str := fmt.Sprintf("%s:%v", *port.Protocol, port.Port.IntVal)
	if port.Port.StrVal != "" {
		str += fmt.Sprintf(" (%s)", port.Port.StrVal)
	}
	return str
}

func fmtPeer(peer v1.NetworkPolicyPeer) string {
	if peer.IPBlock != nil {
		str := peer.IPBlock.CIDR
		if len(peer.IPBlock.Except) > 0 {
			str += fmt.Sprintf(" except %s", strings.Join(peer.IPBlock.Except, ","))
		}
		return str
	}

	str := "all pods"
	return str
}
