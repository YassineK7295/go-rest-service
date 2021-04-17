package main

import (
	"fmt"
	"os"
)

// Entrypoint
func main() {
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Initializes the config and the app, and starts the server
func start() error {
	config, err := InitializeConfig()
	if err != nil {
		return err
	}

	app := App{}
	if err := app.Initialize(config); err != nil {
		return err
	}

	return app.Run(config.SERVE_ADDR)
}
