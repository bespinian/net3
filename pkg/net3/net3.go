package net3

import "k8s.io/client-go/kubernetes"

type Net3 struct {
	k8s kubernetes.Interface
}

// New creates a net3 application.
func New(k8s kubernetes.Interface) *Net3 {
	c := &Net3{
		k8s: k8s,
	}
	return c
}
