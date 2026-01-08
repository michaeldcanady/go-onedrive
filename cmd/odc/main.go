package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal/cmd/root"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ctx := context.Background()

	rootCmd, err := root.CreateRootCmd()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return 1
	}

	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		if isAuthRequired(err) {
			err = errors.New("authentication required. Run `go-onedrive auth login`")
		}
		fmt.Printf("ERROR: %s\n", err)
		return 1
	}

	return 0
}
