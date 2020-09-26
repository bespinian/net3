package net3

import v1 "k8s.io/api/core/v1"

func (n *net3) addLogProxy(podSpec *v1.PodSpec, port, destPort int32) {
	proxyContainer := v1.Container{

		Ports: []v1.ContainerPort{
			{
				ContainerPort: port,
			},
		},
		Env: []v1.EnvVar{
			{
				Name:  "SOURCE_PORT",
				Value: string(port),
			},
			{
				Name:  "DESTINATION_PORT",
				Value: string(destPort),
			},
		},
		Image: "busybox",
		Name:  "net3-log-proxy",
	}

	extendedContainers := append(podSpec.Containers, proxyContainer)
	podSpec.Containers = extendedContainers
}
