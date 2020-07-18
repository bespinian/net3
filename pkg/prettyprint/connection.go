package prettyprint

import (
	"fmt"
	"strings"
)

// Connection prints a connection as a downwards arrow.
func Connection(ports []string) {
	fmt.Println("      │")
	fmt.Println("      │")
	fmt.Printf("      │ %s\n", strings.Join(ports, ", "))
	fmt.Println("      │")
	fmt.Println("      V")
}
