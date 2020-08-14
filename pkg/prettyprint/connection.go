package prettyprint

import (
	"fmt"
	"strings"
)

// Connection prints a connection as a downwards arrow.
func Connection(ports []string, isOpen bool) {
	lines := "      │\n"
	lines += "      │\n"
	lines += fmt.Sprintf("      │ %s\n", strings.Join(ports, ", "))
	lines += "      │\n"
	lines += "      V"

	if isOpen {
		lines = asGreen(lines)
	} else {
		lines = asRed(lines)
	}

	fmt.Println(lines)
}
