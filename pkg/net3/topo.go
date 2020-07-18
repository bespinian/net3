package net3

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Topo lists the topology of a connection path.
func (n *net3) Topo(namespace, src, dest string) error {
	srcPod, err := n.k8s.CoreV1().Pods(namespace).Get(context.Background(), src, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting source pod: %w", err)
	}

	netPols, err := n.k8s.NetworkingV1().NetworkPolicies(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing egress network policies: %w", err)
	}

	fmt.Println("source: ", srcPod.Name)
	fmt.Printf("pols: %+v\n", netPols)

	return nil
}
