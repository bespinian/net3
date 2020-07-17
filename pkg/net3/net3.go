package net3

import "k8s.io/client-go/kubernetes"

type Net3 interface {
	Topo(src, dest string) error
}

type net3 struct {
	client *kubernetes.Clientset
}

func New(client *kubernetes.Clientset) Net3 {
	c := &net3{
		client: client,
	}
	return c
}
