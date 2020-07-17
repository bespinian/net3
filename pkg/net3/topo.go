package net3

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (n *net3) Topo(src, dest string) error {
	srcPod, err := n.client.CoreV1().Pods("").Get(context.TODO(), src, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error listing pods: %w", err)
	}

	fmt.Println("source: ", srcPod.Name)
	return nil
}
