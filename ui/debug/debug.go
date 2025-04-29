package debug

import "fmt"

// Log prints debug information
func Log(msg string) {
	fmt.Println("[DEBUG]:", msg)
}
