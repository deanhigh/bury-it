// Package main is the entry point for bury-it CLI.
package main

import (
	"os"

	"github.com/deanhigh/bury-it/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
