package main

import (
	"github.com/night-slayer18/goforge/cmd"
)

// main is the entry point of the goforge CLI application.
// It does nothing more than execute the root command from the cmd package.
// This lean structure is a best practice for Cobra applications.[7]
func main() {
	cmd.Execute()
}
