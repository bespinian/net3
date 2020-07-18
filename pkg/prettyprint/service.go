package prettyprint

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

// Service prints a service.
func Service(svc *v1.Service) {
	portStrings := make([]string, 0, len(svc.Spec.Ports))
	targetPortStrings := make([]string, 0, len(svc.Spec.Ports))
	for i, p := range svc.Spec.Ports {
		portStrings = append(portStrings, fmt.Sprintf("%s:%v", p.Protocol, p.Port))
		if (p.Name) != "" {
			portStrings[i] += fmt.Sprintf(" (%s)", p.Name)
		}

		targetPortStrings = append(targetPortStrings, fmt.Sprintf("%s:", p.Protocol))
		if p.TargetPort.IntVal > 0 {
			targetPortStrings[i] += fmt.Sprintf("%v", p.TargetPort.IntVal)
		} else {
			targetPortStrings[i] += p.TargetPort.StrVal
		}
	}

	lines := []string{
		"Service",
		fmt.Sprintf("Name: %s", svc.Name),
		fmt.Sprintf("Namespace: %s", svc.Namespace),
		fmt.Sprintf("Ports: %s", strings.Join(portStrings, ", ")),
		fmt.Sprintf("Target Ports: %s", strings.Join(targetPortStrings, ", ")),
	}
	fmt.Print(asBox(lines))
}
