package net3

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (n *Net3) getServicePods(ctx context.Context, namespace, svcName string) ([]corev1.Pod, error) {
	svc, err := n.k8s.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting service %q in namespace %q: %w", svcName, namespace, err)
	}

	pods, err := n.k8s.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing pods in namespace %q: %w", namespace, err)
	}

	svcPods := make([]corev1.Pod, 0, len(pods.Items))
	for i, p := range pods.Items {
		if doesMatchSelector(svc.Spec.Selector, p.Labels) {
			svcPods = append(svcPods, pods.Items[i])
		}
	}

	return svcPods, nil
}
