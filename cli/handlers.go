package cli

import "fmt"

// Close the CLI
func exitHandler(argv []string, stop chan bool) {
	fmt.Println("Attempting to stop CLI...")
	close(stop)
}
