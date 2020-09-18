package net3

import "fmt"

// Topo lists the topology of a connection path.
func (n *net3) Log(namespace, dest string, port int32) error {
	svcPods, err := n.getServicePods(namespace, dest)
	if err != err {
		return fmt.Errorf("error getting pods for service %q in namespace %q: %w", dest, namespace, err)
	}

	for _, p := range svcPods {
		fmt.Println(p.Name)
	}

	return nil
}
