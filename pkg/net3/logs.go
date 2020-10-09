package net3

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Logs redeploys pods with a proxy container which logs all requests to the specified port.
func (n *net3) Logs(namespace, serviceName string, port int32) error {
	// retrieve destination service
	svc, err := n.k8s.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("service with name %q not found in namespace %q: %w", serviceName, namespace, ErrNotFound)
	}

	// look up destination port and highest destination port
	var destPort, maxDestPort int32
	for _, p := range svc.Spec.Ports {
		if p.Port == port {
			destPort = p.TargetPort.IntVal
		}
		if p.TargetPort.IntVal > maxDestPort {
			maxDestPort = p.TargetPort.IntVal
		}
	}

	if destPort == 0 {
		return fmt.Errorf("service %q does not expose port %v: %w", serviceName, port, ErrNotFound)
	}

	// choose a port for the proxy higher than any current destination port
	proxyPort := maxDestPort + 1

	svcPods, err := n.getServicePods(context.Background(), namespace, serviceName)
	if err != nil {
		return fmt.Errorf("error getting pods for service %q in namespace %q: %w", serviceName, namespace, err)
	}

	if len(svcPods) == 0 {
		return fmt.Errorf("service %q has no pods: %w", serviceName, ErrNotFound)
	}

	// find owner of pods
	samplePod := svcPods[0]
	podOwnerRefs := samplePod.GetOwnerReferences()
	if len(podOwnerRefs) == 0 {
		return fmt.Errorf("pod %q has no owner references: %w", samplePod.Name, ErrUnsupported)
	}

	// pod is owned by something
	podOwnerRef := podOwnerRefs[0]
	if podOwnerRef.Kind != "ReplicaSet" {
		return fmt.Errorf("pod %q has unsupported owner kind %q: %w", samplePod.Name, podOwnerRef.Kind, ErrUnsupported)
	}

	// look up the ReplicaSet and check whether it is owned by a Deployment
	replicaSet, err := n.k8s.AppsV1().ReplicaSets(namespace).Get(context.Background(), podOwnerRef.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting replica set %q: %w", podOwnerRef.Name, err)
	}

	// ReplicaSet is owned by something
	replicaSetOwnerRefs := replicaSet.GetOwnerReferences()
	if len(replicaSetOwnerRefs) == 0 {
		return fmt.Errorf("replica set %q has no owner references: %w", replicaSet.Name, ErrUnsupported)
	}

	replicaSetOwnerRef := replicaSetOwnerRefs[0]
	if replicaSetOwnerRef.Kind != "Deployment" {
		return fmt.Errorf("replica set %q has unsupported owner kind %q: %w", replicaSet.Name, replicaSetOwnerRef.Kind, ErrUnsupported)
	}

	// Retrieve the Deployment, add the proxy container to the pod spec and apply
	deployment, err := n.k8s.AppsV1().Deployments(namespace).Get(context.Background(), replicaSetOwnerRef.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting deployment %q: %w", replicaSetOwnerRef.Name, err)
	}

	deployment.Spec.Template.Spec = containerSpecWithProxy(deployment.Spec.Template.Spec, proxyPort, destPort)
	_, err = n.k8s.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating deployment %q in namespace %q with proxy container: %w", deployment.Name, namespace, err)
	}

	// update the service to forward to the proxy port
	svc.Spec = svcSpecWithProxy(svc.Spec, port, proxyPort)
	_, err = n.k8s.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating service %q in namespace %q with proxy port: %w", svc.Name, namespace, err)
	}

	return nil
}
