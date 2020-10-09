package net3

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func svcSpecWithProxy(service v1.ServiceSpec, port, newTargetPort int32) v1.ServiceSpec {
	updatedPorts := make([]v1.ServicePort, 0)
	for _, p := range service.Ports {
		if p.Port == port {
			p.TargetPort = intstr.FromInt(int(newTargetPort))
		}
		updatedPorts = append(updatedPorts, p)
	}

	service.Ports = updatedPorts

	return service
}
