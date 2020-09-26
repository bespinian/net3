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
		return fmt.Errorf("service %q does not expose port %q: %w", dest, port, ErrNotFound)
	}

	// choose a port for the proxy higher than any current destination port
	proxyPort := maxDestPort + 1

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
					n.addLogProxy(&deployment.Spec.Template.Spec, proxyPort, destPort)
					_, err = n.k8s.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
					if err != nil {
						return fmt.Errorf("cannot update deployment %q in namespace %q with proxy container: %w", deployment.Name, namespace, err)
					}
					// update the service to forward to the proxy port
					n.updateServicePort(&svc.Spec, port, proxyPort)
					_, err = n.k8s.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
					if err != nil {
						return fmt.Errorf("cannot update service %q in namespace %q with proxy container port: %w", svc.Name, namespace, err)
					}
					fmt.Println(deployment.Name)
				}
			}
		}
	}

	for _, p := range svcPods {
		fmt.Println(p.Name)
	}

	return nil
}
