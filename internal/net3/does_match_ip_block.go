package net3

import (
	"fmt"
	"net"

	v1 "k8s.io/api/networking/v1"
)

func doesMatchIPBlock(block v1.IPBlock, address string) (bool, error) {
	_, ipv4Net, err := net.ParseCIDR(block.CIDR)
	if err != nil {
		return false, fmt.Errorf("error parsing CIDR: %w", err)
	}
	if !ipv4Net.Contains([]byte(address)) {
		return false, nil
	}

	for _, e := range block.Except {
		_, eNet, err := net.ParseCIDR(e)
		if err != nil {
			return false, fmt.Errorf("error parsing exception: %w", err)
		}
		if eNet.Contains([]byte(address)) {
			return false, nil
		}
	}

	return true, nil
}
