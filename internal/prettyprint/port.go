package prettyprint

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func portString(protocol v1.Protocol, port intstr.IntOrString) string {
	if port.IntVal == 0 {
		if port.StrVal == "" {
			return "undefined"
		}
		return port.StrVal
	}

	str := fmt.Sprintf("%s:%v", protocol, port.IntVal)
	if port.StrVal != "" {
		str += fmt.Sprintf(" (%s)", port.StrVal)
	}
	return str
}
