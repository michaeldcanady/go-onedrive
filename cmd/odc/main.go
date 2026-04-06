// Package main is the entry point for the odc2 application.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/root"
)

func main() {
	// Setup context with signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Initialize dependency injection container
	container, err := di.NewDefaultContainer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize services: %v\n", err)
		os.Exit(1)
	}

	// Create and execute the root command
	cmd, err := root.CreateRootCmd(container)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.ExecuteContext(ctx); err != nil {
		cliErr := errors.MapToCLI(err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", cliErr)
		os.Exit(cliErr.ExitCode)
	}
}
