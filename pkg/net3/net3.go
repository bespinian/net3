package net3

import "k8s.io/client-go/kubernetes"

// Net3 represents a net3 application.
type Net3 interface {
	Topo(namespace, src, dest string) error
	Logs(namespace, dest string, port int32) error
}

type net3 struct {
	k8s *kubernetes.Clientset
}

// New creates a net3 application.
func New(k8s *kubernetes.Clientset) Net3 {
	c := &net3{
		k8s: k8s,
	}
	return c
}
