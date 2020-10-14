package net3

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const svcAnnotationNamePrefix = "net3.bespinian.io/proxy-container-name"

func svcWithProxy(svc v1.Service, containerName string, port, newTargetPort int32) v1.Service {
	svc.Spec = svcSpecWithTargetPort(svc.Spec, port, newTargetPort)
	annotationName := fmt.Sprintf("%s-%v", svcAnnotationNamePrefix, port)
	svc.Annotations[annotationName] = containerName
	return svc
}

func svcWithoutProxy(svc v1.Service, port, originalTargetPort int32) v1.Service {
	svc.Spec = svcSpecWithTargetPort(svc.Spec, port, originalTargetPort)
	annotationName := fmt.Sprintf("%s-%v", svcAnnotationNamePrefix, port)
	delete(svc.Annotations, annotationName)
	return svc
}

func svcSpecWithTargetPort(service v1.ServiceSpec, port, newTargetPort int32) v1.ServiceSpec {
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
