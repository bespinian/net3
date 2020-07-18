package prettyprint

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

// Pod prints a pod.
func Pod(pod *v1.Pod) {
	var portStrings []string
	for _, c := range pod.Spec.Containers {
		for _, p := range c.Ports {
			str := fmt.Sprintf("%s:%v", p.Protocol, p.ContainerPort)
			if p.Name != "" {
				str += fmt.Sprintf(" (%s)", p.Name)
			}
			portStrings = append(portStrings, str)
		}
	}
	lines := []string{
		"Pod",
		fmt.Sprintf("Name: %s", pod.Name),
		fmt.Sprintf("Namespace: %s", pod.Namespace),
		fmt.Sprintf("Status: %s", pod.Status.Phase),
	}
	if len(portStrings) > 0 {
		lines = append(lines, fmt.Sprintf("Ports: %s", strings.Join(portStrings, ", ")))
	}
	fmt.Print(asBox(lines))
}
