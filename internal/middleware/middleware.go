package middleware

import "github.com/spf13/cobra"

type CobraMiddleware = func(*cobra.Command)
