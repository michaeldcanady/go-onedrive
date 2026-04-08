package middleware

import "github.com/spf13/cobra"

// ApplyMiddlewareRecursively applies the given middleware functions to the command and all of its subcommands recursively.
func ApplyMiddlewareRecursively(cmd *cobra.Command, mws ...CobraMiddleware) {

	for _, mw := range mws {
		mw(cmd)
	}
	for _, c := range cmd.Commands() {
		ApplyMiddlewareRecursively(c, mws...)
	}
}
