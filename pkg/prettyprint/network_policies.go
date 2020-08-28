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
			fmt.Sprintf("Name:       %s", p.Name),
			fmt.Sprintf("Namespace:  %s", p.Namespace),
			fmt.Sprintf("Allowing:   %v", isAllowed),
		}

		if policyType == v1.PolicyTypeIngress {
			for _, r := range p.Spec.Ingress {
				lines = append(lines, fmt.Sprintf("Rule:       %s", fmtIngressRule(r)))
			}
		}

		fmt.Print(asBox(lines))
	}
}

func fmtIngressRule(rule v1.NetworkPolicyIngressRule) string {
	str := "Allow "

	portStrings := make([]string, 0, len(rule.Ports))
	for _, p := range rule.Ports {
		portStrings = append(portStrings, fmtPort(p))
	}
	if len(portStrings) == 0 {
		portStrings = append(portStrings, "any traffic")
	}
	str += strings.Join(portStrings, ",")

	fromStrings := make([]string, 0, len(rule.From))
	for _, f := range rule.From {
		fromStrings = append(fromStrings, fmtPeer(f))
	}
	str += fmt.Sprintf(" from %s", strings.Join(fromStrings, ","))
	return str

}

func fmtPort(port v1.NetworkPolicyPort) string {
	if port.Protocol == nil && port.Port == nil {
		return "any TCP traffic"
	}
	if port.Port == nil {
		return fmt.Sprintf("any %s traffic", *port.Protocol)
	}
	if port.Port.IntVal == 0 && port.Port.StrVal == "" {
		return fmt.Sprintf("any %s traffic", *port.Protocol)
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
	str := ""

	if peer.PodSelector == nil {
		str = "any pod"
	} else {
		labelStrings := make([]string, 0, len(peer.PodSelector.MatchLabels))
		for k, v := range peer.PodSelector.MatchLabels {
			labelStrings = append(labelStrings, fmt.Sprintf("%s=%s", k, v))
		}
		str += fmt.Sprintf("pods [%s]", strings.Join(labelStrings, ","))
	}

	if peer.IPBlock != nil {
		str := peer.IPBlock.CIDR
		if len(peer.IPBlock.Except) > 0 {
			str += fmt.Sprintf(" except %s", strings.Join(peer.IPBlock.Except, ","))
		}
		return str
	}

	return str
}
