package net3

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
)

const (
	proxyContainerName    = "net3-log-proxy"
	proxyEnvVarPort       = "NET3_HTTP_PROXY_PORT"
	proxyEnvVarTargetPort = "NET3_HTTP_PROXY_TARGET_PORT"
)

func podSpecWithProxy(podSpec v1.PodSpec, proxyPort, targetPort int32) v1.PodSpec {
	proxyContainer := v1.Container{
		Name:            proxyContainerName,
		Image:           "bespinian/net3-http-proxy:0.0.1",
		ImagePullPolicy: v1.PullAlways,
		Ports: []v1.ContainerPort{
			{ContainerPort: proxyPort},
		},
		Env: []v1.EnvVar{
			{
				Name:  proxyEnvVarPort,
				Value: strconv.Itoa(int(proxyPort)),
			},
			{
				Name:  proxyEnvVarTargetPort,
				Value: strconv.Itoa(int(targetPort)),
			},
		},
	}

	podSpec.Containers = append(podSpec.Containers, proxyContainer)

	return podSpec
}

func podSpecWithoutProxy(podSpec v1.PodSpec) (v1.PodSpec, int32, error) {
	var originalPort int32
	containers := make([]v1.Container, 0)

	for _, c := range podSpec.Containers {
		if c.Name == proxyContainerName {
			for _, e := range c.Env {
				if e.Name == proxyEnvVarTargetPort {
					v, err := strconv.Atoi(e.Value)
					if err != nil {
						err = fmt.Errorf("error converting target port to int: %w", err)
						return v1.PodSpec{}, 0, err
					}
					originalPort = int32(v)
				}
			}
		} else {
			containers = append(containers, c)
		}
	}
	if len(containers) == len(podSpec.Containers) {
		err := fmt.Errorf("could not find proxy container in pod spec: %w", ErrNotFound)
		return v1.PodSpec{}, 0, err
	}
	if originalPort == 0 {
		err := fmt.Errorf("could not find environment variable %s in proxy container spec: %w", proxyEnvVarTargetPort, ErrNotFound)
		return v1.PodSpec{}, 0, err
	}

	podSpec.Containers = containers

	return podSpec, originalPort, nil
}
