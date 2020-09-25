package net3

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Log redeploys pods with a proxy container which logs all requests to the specified port.
func (n *net3) Log(namespace, dest string, port int32) error {
	// retrieve destination service
	svc, err := n.k8s.CoreV1().Services(namespace).Get(context.Background(), dest, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("service with name %q not found in namespace %q: %w", dest, namespace, ErrNotFound)
	}

	// look up destination port
	var destPort int32
	for _, p := range svc.Spec.Ports {
		if p.Port == port {
			destPort = p.TargetPort.IntVal
			break
		}
	}

	if destPort == 0 {
		return fmt.Errorf("service %q does not expose port %q: %w", dest, port, ErrNotFound)
	}

	svcPods, err := n.getServicePods(namespace, dest)
	if err != nil {
		return fmt.Errorf("error getting pods for service %q in namespace %q: %w", dest, namespace, err)
	}

	if len(svcPods) == 0 {
		return fmt.Errorf("service %q has no pods: %w", dest, ErrNotFound)
	}

	samplePod := svcPods[0]

	// find owner of pods
	if len(samplePod.GetOwnerReferences()) > 0 {
		// pod is owned by something
		podOwnerRef := samplePod.GetOwnerReferences()[0]
		if podOwnerRef.Kind == "ReplicaSet" {
			// look up the ReplicaSet and check whether it is owned by a Deployment
			replicaSet, err := n.k8s.AppsV1().ReplicaSets(namespace).Get(context.Background(), podOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("no ReplicaSet %q matching owner ref found: %w", podOwnerRef.Name, err)
			}
			if len(replicaSet.GetOwnerReferences()) > 0 {
				// ReplicaSet is owned by something
				replicaSetOwnerRef := replicaSet.GetOwnerReferences()[0]
				if replicaSetOwnerRef.Kind == "Deployment" {
					// Retrieve the Deployment, add the proxy container to the pod spec and apply
					deployment, err := n.k8s.AppsV1().Deployments(namespace).Get(context.Background(), replicaSetOwnerRef.Name, metav1.GetOptions{})
					if err != nil {
						return fmt.Errorf("no Deployment %q matching owner ref found: %w", replicaSetOwnerRef.Name, err)
					}
					fmt.Println(deployment.Name)
				} else {
					// any other possibilities (don't know any)?
				}
			} else {
				// ReplicaSet has been created manually. Add the proxy container to the pod spec and apply
			}
		} else {
			// any other possibilities (don't know any)?
		}
	} else {
		// pod has been created manually, conclude that all others have too and replace them one by one
	}

	for _, p := range svcPods {
		fmt.Println(p.Name)
	}

	return nil
}
