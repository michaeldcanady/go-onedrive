package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/michaeldcanady/go-onedrive/internal2/app/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/root"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ctx := context.Background()

	container := di.NewContainer()

	rootCmd, err := root.CreateRootCmd(container)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return 1
	}

	if cmd, err := rootCmd.ExecuteContextC(ctx); err != nil {
		if isAuthRequired(err) {
			err = errors.New("authentication required. Run `odc auth login`")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", cmd.CalledAs(), err)
		return 1
	}

	return 0
}
