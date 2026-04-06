package middleware

import "github.com/spf13/cobra"

// ApplyMiddlewareRecursively applies the given [CobraMiddleware] to the provided [cobra.Command] and all of its subcommands.
func ApplyMiddlewareRecursively(cmd *cobra.Command, mws ...CobraMiddleware) {

	for _, mw := range mws {
		mw(cmd)
	}
	for _, c := range cmd.Commands() {
		ApplyMiddlewareRecursively(c, mws...)
	}
}
