package main

import (
	"fmt"
	"github.com/erniebilling/gcp-orphans/cmd"
	"os"
)

func main() {
	root := cmd.CreateRootCommand()
	fmt.Println() // Print a blank line before output for readability

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
	fmt.Println() // Print a blank line after output for readability
}
