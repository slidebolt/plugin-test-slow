package main

import (
	"log"

	"github.com/slidebolt/plugin-test-slow/pkg"
	runner "github.com/slidebolt/sdk-runner"
)

func main() {
	if err := runner.RunCLI(func() runner.Plugin {
		// Create the core provider (decoupled from SDK).
		provider := pkg.NewSlowProvider()
		// Create the adapter that bridges the provider to the SDK.
		return NewSlowPluginAdapter(provider)
	}); err != nil {
		log.Fatal(err)
	}
}
