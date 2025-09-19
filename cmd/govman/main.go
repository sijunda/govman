// cmd/govman/main.go
package main

import (
	"fmt"
	"os"

	_cli "github.com/sijunda/govman/internal/cli"
)

func main() {
	if err := _cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
