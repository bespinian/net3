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
	// Source
	srcPod, err := n.k8s.CoreV1().Pods(namespace).Get(context.Background(), src, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting source pod in namespace %q: %w", namespace, err)
	}
	srcNetPolList, err := n.k8s.NetworkingV1().NetworkPolicies(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing source network policies in namespace %q: %w", namespace, err)
	}

	// Destination
	destination, err := NewDestination(dest, namespace)
	if err != nil {
		return fmt.Errorf("error parsing destination: %w", err)
	}
	if destination.Kind != DestinationKindService {
		return fmt.Errorf("destination kind %s currently not supported", destination.Kind)
	}
	svc, err := n.k8s.CoreV1().Services(destination.Namespace).Get(context.Background(), destination.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting destination service in namespace %q: %w", destination.Namespace, err)
	}
	destNetPolList, err := n.k8s.NetworkingV1().NetworkPolicies(destination.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing destination network policies in namespace %q: %w", destination.Namespace, err)
	}
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

	// Egress connection from source
	egressPols := make([]networkingv1.NetworkPolicy, 0, len(srcNetPolList.Items))
	for _, p := range srcNetPolList.Items {
		if doesMatchSelector(p.Spec.PodSelector.MatchLabels, srcPod.Labels) {
			if len(p.Spec.Egress) > 0 {
				egressPols = append(egressPols, p)
			}
		}
	}

	// Egress connection from service
	svcTargetPort := 0
	for _, p := range svc.Spec.Ports {
		if destination.Port == int(p.Port) {
			svcTargetPort = p.TargetPort.IntValue()
		}
	}
	doesSvcPortExist := svcTargetPort != 0

	// Ingress connection to destination
	ingressPols := make([]networkingv1.NetworkPolicy, 0, len(destNetPolList.Items))
	for _, p := range destNetPolList.Items {
		if doesMatchSelector(p.Spec.PodSelector.MatchLabels, destPod.Labels) {
			if len(p.Spec.Ingress) > 0 {
				ingressPols = append(ingressPols, p)
			}
		}
	}
	allowingPolicies := make([]networkingv1.NetworkPolicy, 0, len(ingressPols))
	denyingPolicies := make([]networkingv1.NetworkPolicy, 0, len(ingressPols))
	for _, p := range ingressPols {
		doesPolMatch := false
		for _, r := range p.Spec.Ingress {
			doesRuleMatch, err := n.doesMatchIngressRule(r, srcPod, svcTargetPort)
			if err != nil {
				return fmt.Errorf("error checking if ingress rule matches: %w", err)
			}
			if doesRuleMatch {
				doesPolMatch = true
				break
			}
		}
		if doesPolMatch {
			allowingPolicies = append(allowingPolicies, p)
		} else {
			denyingPolicies = append(denyingPolicies, p)
		}
	}

	fmt.Println("")
	prettyprint.Pod(srcPod)
	if len(egressPols) > 0 {
		prettyprint.Connection([]string{destination.FullPort()}, true)
		prettyprint.NetworkPolicies(networkingv1.PolicyTypeEgress, egressPols, true)
	}
	prettyprint.Connection([]string{destination.FullPort()}, doesSvcPortExist)
	prettyprint.Service(svc)
	prettyprint.Connection([]string{fmt.Sprintf("%s:%v", "TCP", svcTargetPort)}, true)
	if len(ingressPols) > 0 {
		if len(allowingPolicies) > 0 {
			prettyprint.NetworkPolicies(networkingv1.PolicyTypeIngress, allowingPolicies, true)
		} else {
			prettyprint.NetworkPolicies(networkingv1.PolicyTypeIngress, denyingPolicies, false)
		}
		prettyprint.Connection([]string{fmt.Sprintf("%s:%v", "TCP", svcTargetPort)}, true)
	}
	prettyprint.Pod(destPod)
	fmt.Println("")

	return nil
}
