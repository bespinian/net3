package net3

import (
	"fmt"
	"strconv"
	"strings"
)

const defaultHTTPPort = 80

// Destination is a destination for a request.
type Destination struct {
	Kind      DestinationKind
	Name      string
	Namespace string
	Protocol  string
	Domain    string
	Port      int32
}

// DestinationKind is a type of destination.
type DestinationKind string

const (
	DestinationKindService DestinationKind = "service"
)

// NewDestination creates a new destination from an address.
func NewDestination(address, defaultNamespace string) (*Destination, error) {
	d := Destination{
		Kind:      DestinationKindService,
		Name:      address,
		Namespace: defaultNamespace,
		Protocol:  "TCP",
		Domain:    address,
		Port:      defaultHTTPPort,
	}

	addressParts := strings.Split(address, ":")
	if len(addressParts) > 1 {
		d.Name = addressParts[0]
		d.Domain = addressParts[0]
		port, err := strconv.Atoi(addressParts[len(addressParts)-1]) //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		d.Port = int32(port)
	}

	domainParts := strings.Split(d.Domain, ".")
	if len(domainParts) > 1 {
		d.Name = domainParts[0]
		d.Namespace = domainParts[1]
	}

	return &d, nil
}

// FullPort returns the protocol and port of the destionation.
func (d *Destination) FullPort() string {
	return fmt.Sprintf("%s:%v", d.Protocol, d.Port)
}
