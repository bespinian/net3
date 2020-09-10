package prettyprint

import (
	"fmt"
	"strings"
)

const arrowLine = "      │\n"

// Connection prints a connection as a downwards arrow.
func Connection(ports []string, isOpen bool) {
	lines := arrowLine
	lines += arrowLine
	lines += fmt.Sprintf("      │ %s\n", strings.Join(ports, ", "))
	lines += arrowLine
	lines += "      V"

	if isOpen {
		lines = asGreen(lines)
	} else {
		lines = asRed(lines)
	}

	fmt.Println(lines)
}
