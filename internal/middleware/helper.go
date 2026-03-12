package middleware

import "github.com/spf13/cobra"

func ApplyMiddlewareRecursively(cmd *cobra.Command, mws ...CobraMiddleware) {

	for _, mw := range mws {
		mw(cmd)
	}
	for _, c := range cmd.Commands() {
		ApplyMiddlewareRecursively(c, mws...)
	}
}
