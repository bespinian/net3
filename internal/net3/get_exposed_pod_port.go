package net3

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getExposedPodPort(pod *v1.Pod, port intstr.IntOrString) (int32, bool) {
	for _, c := range pod.Spec.Containers {
		for _, p := range c.Ports {
			if port.IntVal != 0 {
				if p.ContainerPort == port.IntVal {
					return p.ContainerPort, true
				}
			}
			if port.StrVal != "" {
				if p.Name == port.StrVal {
					return p.ContainerPort, true
				}
			}
		}
	}
	return 0, false
}
