package net3

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// AddProxy redeploys pods with a proxy container which logs all requests to the specified port.
func (n *Net3) AddProxy(namespace, serviceName, containerName, image string, port int32) error {
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

	// look up target port and highest target port
	var (
		originalTargetPort        intstr.IntOrString
		targetPort, maxTargetPort int32
		ok                        bool
	)
	for _, p := range svc.Spec.Ports {
		if p.Port == port {
			fmt.Printf("%+v\n", p.TargetPort)
			originalTargetPort = p.TargetPort
			targetPort, ok = getExposedPodPort(&samplePod, p.TargetPort)
		}
		if targetPort > maxTargetPort {
			maxTargetPort = targetPort
		}
	}

	if !ok {
		return fmt.Errorf("service %q does not expose port %v: %w", serviceName, port, ErrNotFound)
	}

	// choose a port for the proxy higher than any current target port
	proxyPort := maxTargetPort + 1

	// update the service to forward to the proxy port
	*svc = svcWithProxy(*svc, containerName, port, originalTargetPort, intstr.FromInt(int(proxyPort)))
	_, err = n.k8s.CoreV1().Services(namespace).Update(context.Background(), svc, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating service %q in namespace %q with proxy port: %w", svc.Name, namespace, err)
	}

	// update deployment pod spec
	deployment.Spec.Template.Spec = podSpecWithProxy(deployment.Spec.Template.Spec, containerName, image, proxyPort, targetPort)
	_, err = n.k8s.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating deployment %q in namespace %q with proxy container: %w", deployment.Name, namespace, err)
	}

	labelStrings := make([]string, 0, len(svc.Spec.Selector))
	for k, v := range svc.Spec.Selector {
		labelStrings = append(labelStrings, fmt.Sprintf("%s=%s", k, v))
	}
	logCommand := fmt.Sprintf("kubectl -n %s logs -l %s -c %s -f", namespace, strings.Join(labelStrings, ","), containerName)

	fmt.Printf("Adding log proxy to pods of service %q as container %q\n", svc.Name, containerName)
	fmt.Printf("Get the request and response logs with %q\n", logCommand)

	return nil
}
