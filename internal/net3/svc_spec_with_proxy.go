package net3

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	svcContainerNameAnnotationPrefix      = "net3.bespinian.io/proxy-container-name"
	svcOriginalTargetPortAnnotationPrefix = "net3.bespinian.io/original-target-port"
)

func svcWithProxy(svc v1.Service, containerName string, port int32, originalTargetPort, newTargetPort intstr.IntOrString) v1.Service {
	svc.Spec = svcSpecWithTargetPort(svc.Spec, port, newTargetPort)

	originalTargetPortStr := originalTargetPort.StrVal
	if originalTargetPortStr == "" {
		originalTargetPortStr = fmt.Sprintf("%v", originalTargetPort.IntVal)
	}

	containerNameAnnotation := fmt.Sprintf("%s-%v", svcContainerNameAnnotationPrefix, port)
	originalTargetPortAnnotation := fmt.Sprintf("%s-%v", svcOriginalTargetPortAnnotationPrefix, port)
	svc.Annotations[containerNameAnnotation] = containerName
	svc.Annotations[originalTargetPortAnnotation] = originalTargetPortStr

	return svc
}

func svcWithoutProxy(svc v1.Service, port int32) v1.Service {
	containerNameAnnotation := fmt.Sprintf("%s-%v", svcContainerNameAnnotationPrefix, port)
	originalTargetPortAnnotation := fmt.Sprintf("%s-%v", svcOriginalTargetPortAnnotationPrefix, port)

	var originalTargetPort intstr.IntOrString
	targetPortInt, err := strconv.Atoi(svc.Annotations[originalTargetPortAnnotation])
	if err == nil {
		originalTargetPort = intstr.FromInt(targetPortInt)
	} else {
		originalTargetPort = intstr.FromString(svc.Annotations[originalTargetPortAnnotation])
	}

	svc.Spec = svcSpecWithTargetPort(svc.Spec, port, originalTargetPort)

	delete(svc.Annotations, containerNameAnnotation)
	delete(svc.Annotations, originalTargetPortAnnotation)

	return svc
}

func svcSpecWithTargetPort(service v1.ServiceSpec, port int32, newTargetPort intstr.IntOrString) v1.ServiceSpec {
	updatedPorts := make([]v1.ServicePort, 0)
	for _, p := range service.Ports {
		if p.Port == port {
			p.TargetPort = newTargetPort
		}
		updatedPorts = append(updatedPorts, p)
	}

	service.Ports = updatedPorts

	return service
}
