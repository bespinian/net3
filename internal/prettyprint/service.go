package prettyprint

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// Service prints a service.
func Service(svc *v1.Service) {
	lines := []string{
		"Service",
		fmt.Sprintf("Name:          %s", svc.Name),
		fmt.Sprintf("Namespace:     %s", svc.Namespace),
	}

	if svc.Spec.ClusterIP != "" {
		lines = append(lines, fmt.Sprintf("IP:            %s", svc.Spec.ClusterIP))
	}

	for _, p := range svc.Spec.Ports {
		str := fmt.Sprintf("%s:%v", p.Protocol, p.Port)
		if (p.Name) != "" {
			str += fmt.Sprintf(" (%s)", p.Name)
		}
		str += fmt.Sprintf(" -> %s", portString(p.Protocol, p.TargetPort))
		lines = append(lines, fmt.Sprintf("Port:          %s", str))
	}

	fmt.Print(asBox(lines))
}
