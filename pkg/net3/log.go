package net3

import (
	"fmt"
)

// Log redeploys pods with a proxy container which logs all requests to the specified port.
func (n *net3) Log(namespace, dest string, port int32) error {
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
		ownerRef := samplePod.GetOwnerReferences()[0]
		if ownerRef.Kind == "ReplicaSet" {
			// look up the ReplicaSet and check whether it is owned by a Deployment

		} else {
			// any other possibilities?
		}
	} else {
		// pod has been created manually, conclude that all others have too and replace them one by one
	}

	for _, p := range svcPods {
		fmt.Println(p.Name)
	}

	return nil
}
