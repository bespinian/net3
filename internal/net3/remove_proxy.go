package net3

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RemoveProxy removes an existing net3 proxy from a service.
func (n *Net3) RemoveProxy(namespace, serviceName string, port int32) error {
	// retrieve destination service
	svc, err := n.k8s.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("service with name %q not found in namespace %q: %w", serviceName, namespace, ErrNotFound)
	}

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

	annotationName := fmt.Sprintf("%s-%v", svcContainerNameAnnotationPrefix, port)
	containerName, ok := svc.Annotations[annotationName]
	if !ok {
		return fmt.Errorf("error getting container name from service annotation %q: %w", annotationName, ErrNotFound)
	}

	// update the service to forward to the original target port
	*svc = svcWithoutProxy(*svc, port)
	_, err = n.k8s.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating service %q in namespace %q: %w", svc.Name, namespace, err)
	}

	// update deployment pod spec
	podSpec, err := podSpecWithoutProxy(deployment.Spec.Template.Spec, containerName)
	if err != nil {
		return fmt.Errorf("error removing proxy container from pods of deployment %q: %w", deployment.Name, err)
	}
	deployment.Spec.Template.Spec = podSpec
	_, err = n.k8s.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating deployment %q in namespace %q: %w", deployment.Name, namespace, err)
	}

	fmt.Printf("Removing log proxy from pods of service %q\n", svc.Name)

	return nil
}
