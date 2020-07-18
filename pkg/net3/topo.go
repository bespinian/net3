package net3

import (
	"context"
	"fmt"

	"github.com/bespinian/net3/pkg/prettyprint"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Topo lists the topology of a connection path.
func (n *net3) Topo(namespace, src, dest string) error {
	fmt.Println("")

	// Source pod
	srcPod, err := n.k8s.CoreV1().Pods(namespace).Get(context.Background(), src, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting source pod in namespace %q: %w", namespace, err)
	}
	prettyprint.Pod(srcPod)

	// Egress connection
	destination, err := NewDestination(dest, namespace)
	if err != nil {
		return fmt.Errorf("error parsing destination: %w", err)
	}
	if destination.Kind != DestinationKindService {
		return fmt.Errorf("destination kind %s currently not supported", destination.Kind)
	}
	prettyprint.Connection([]string{destination.FullPort()})

	// Egress network policies
	srcNetPolList, err := n.k8s.NetworkingV1().NetworkPolicies(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing source network policies in namespace %q: %w", namespace, err)
	}
	egressPols := make([]networkingv1.NetworkPolicy, 0, len(srcNetPolList.Items))
	for _, p := range srcNetPolList.Items {
		if doesMatchSelector(p.Spec.PodSelector.MatchLabels, srcPod.Labels) {
			if len(p.Spec.Egress) > 0 {
				egressPols = append(egressPols, p)
			}
		}
	}
	if len(egressPols) > 0 {
		prettyprint.NetworkPolicies(networkingv1.PolicyTypeEgress, egressPols)
		prettyprint.Connection([]string{destination.FullPort()})
	}

	// Destination service
	svc, err := n.k8s.CoreV1().Services(destination.Namespace).Get(context.Background(), destination.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting destination service in namespace %q: %w", destination.Namespace, err)
	}
	prettyprint.Service(svc)

	// Egress connection from service
	targetPorts := make([]string, 0, len(svc.Spec.Ports))
	for i, p := range svc.Spec.Ports {
		if p.TargetPort.IntVal > 0 {
			targetPorts = append(targetPorts, fmt.Sprintf("%s:%v", p.Protocol, p.TargetPort.IntVal))
			if p.TargetPort.StrVal != "" {
				targetPorts[i] += fmt.Sprintf(" (%s)", p.TargetPort.StrVal)
			}
		} else {
			targetPorts = append(targetPorts, p.TargetPort.StrVal)
		}
	}
	prettyprint.Connection(targetPorts)

	// Destination pod
	destPods, err := n.k8s.CoreV1().Pods(destination.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing destination pods in namespace %q: %w", destination.Namespace, err)
	}
	var destPod *corev1.Pod
	for _, p := range destPods.Items {
		if doesMatchSelector(svc.Spec.Selector, p.Labels) {
			destPod = &p
			break
		}
	}
	if destPod == nil {
		return fmt.Errorf("no matching pod found for service %q in namespace %q", svc.Name, svc.Namespace)
	}

	// Ingress network policies
	destNetPolList, err := n.k8s.NetworkingV1().NetworkPolicies(destination.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing destination network policies in namespace %q: %w", destination.Namespace, err)
	}
	ingressPols := make([]networkingv1.NetworkPolicy, 0, len(destNetPolList.Items))
	for _, p := range destNetPolList.Items {
		if doesMatchSelector(p.Spec.PodSelector.MatchLabels, destPod.Labels) {
			if len(p.Spec.Ingress) > 0 {
				ingressPols = append(ingressPols, p)
			}
		}
	}
	if len(ingressPols) > 0 {
		prettyprint.NetworkPolicies(networkingv1.PolicyTypeIngress, ingressPols)
		prettyprint.Connection(targetPorts)
	}

	prettyprint.Pod(destPod)

	fmt.Println("")
	return nil
}
