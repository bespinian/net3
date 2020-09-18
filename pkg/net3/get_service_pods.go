package net3

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (n *net3) getServicePods(svcName, namespace string) ([]corev1.Pod, error) {
	svc, err := n.k8s.CoreV1().Services(namespace).Get(context.Background(), svcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting destination service in namespace %q: %w", namespace, err)
	}
	destPods, err := n.k8s.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing destination pods in namespace %q: %w", namespace, err)
	}
	svcPods := make([]corev1.Pod, 0)
	for i, p := range destPods.Items {
		if doesMatchSelector(svc.Spec.Selector, p.Labels) {
			svcPods = append(svcPods, destPods.Items[i])
		}
	}
	return svcPods, nil
}
