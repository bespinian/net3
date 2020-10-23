package net3

import (
	"fmt"
	"strconv"

	v1 "k8s.io/api/core/v1"
)

const (
	proxyEnvVarPort       = "NET3_HTTP_PROXY_PORT"
	proxyEnvVarTargetPort = "NET3_HTTP_PROXY_TARGET_PORT"
)

func podSpecWithProxy(podSpec v1.PodSpec, containerName, image string, proxyPort, targetPort int32) v1.PodSpec {
	proxyContainer := v1.Container{
		Name:            containerName,
		Image:           image,
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

func podSpecWithoutProxy(podSpec v1.PodSpec, containerName string) (v1.PodSpec, error) {
	var originalPort int
	containers := make([]v1.Container, 0)

	for _, c := range podSpec.Containers {
		if c.Name == containerName {
			for _, e := range c.Env {
				if e.Name == proxyEnvVarTargetPort {
					v, err := strconv.Atoi(e.Value)
					if err != nil {
						err = fmt.Errorf("error converting target port to int: %w", err)
						return v1.PodSpec{}, err
					}
					originalPort = v
				}
			}
		} else {
			containers = append(containers, c)
		}
	}
	if len(containers) == len(podSpec.Containers) {
		err := fmt.Errorf("could not find proxy container in pod spec: %w", ErrNotFound)
		return v1.PodSpec{}, err
	}
	if originalPort == 0 {
		err := fmt.Errorf("could not find environment variable %s in proxy container spec: %w", proxyEnvVarTargetPort, ErrNotFound)
		return v1.PodSpec{}, err
	}

	podSpec.Containers = containers

	return podSpec, nil
}
