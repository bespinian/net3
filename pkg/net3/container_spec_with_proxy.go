package net3

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
)

func containerSpecWithProxy(podSpec v1.PodSpec, proxyPort, targetPort int32) v1.PodSpec {
	proxyContainer := v1.Container{
		Name:            "net3-log-proxy",
		Image:           "bespinian/net3-http-proxy:0.0.1",
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{
			{ContainerPort: proxyPort},
		},
		Env: []v1.EnvVar{
			{
				Name:  "NET3_HTTP_PROXY_PORT",
				Value: strconv.Itoa(int(proxyPort)),
			},
			{
				Name:  "NET3_HTTP_PROXY_TARGET_PORT",
				Value: strconv.Itoa(int(targetPort)),
			},
		},
	}

	podSpec.Containers = append(podSpec.Containers, proxyContainer)

	return podSpec
}
