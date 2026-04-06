package middleware

import "github.com/spf13/cobra"

// CobraMiddleware defines the type for middleware functions that can be applied to [cobra.Command].
type CobraMiddleware = func(*cobra.Command)
