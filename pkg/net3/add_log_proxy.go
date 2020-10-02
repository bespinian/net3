package net3

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
)

func (n *net3) addLogProxy(podSpec *v1.PodSpec, port, destPort int32) {
	proxyContainer := v1.Container{

		Ports: []v1.ContainerPort{
			{
				ContainerPort: port,
			},
		},
		Env: []v1.EnvVar{
			{
				Name:  "NET3_HTTP_PROXY_PORT",
				Value: strconv.Itoa(int(port)),
			},
			{
				Name:  "NET3_HTTP_PROXY_TARGET_PORT",
				Value: strconv.Itoa(int(destPort)),
			},
		},
		Image: "bespinian/net3-http-proxy:0.0.1",
		Name:  "net3-log-proxy",
	}

	extendedContainers := append(podSpec.Containers, proxyContainer)
	podSpec.Containers = extendedContainers
}
