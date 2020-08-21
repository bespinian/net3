package prettyprint

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// Pod prints a pod.
func Pod(pod *v1.Pod) {
	lines := []string{
		"Pod",
		fmt.Sprintf("Name:       %s", pod.Name),
		fmt.Sprintf("Namespace:  %s", pod.Namespace),
		fmt.Sprintf("Status:     %s", pod.Status.Phase),
		fmt.Sprintf("IP:         %s", pod.Status.PodIP),
	}

	for _, c := range pod.Spec.Containers {
		for _, p := range c.Ports {
			str := fmt.Sprintf("%s:%v", p.Protocol, p.ContainerPort)
			if p.Name != "" {
				str += fmt.Sprintf(" (%s)", p.Name)
			}
			lines = append(lines, fmt.Sprintf("Port:       %s", str))
		}
	}

	fmt.Print(asBox(lines))
}
