package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// Handler defines the lifecycle of a CLI command.
type Handler[T any] interface {
	Validate(ctx *T) error
	Execute(ctx *T) error
	Finalize(ctx *T) error
}

// CommandConfig holds the configuration for creating a new command.
type CommandConfig[T any] struct {
	Use               string
	Short             string
	Long              string
	Args              cobra.PositionalArgs
	ValidArgsFunction func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
	Handler           Handler[T]
	Options           *T
	CtxFunc           func(context.Context, *T) *T
}

// NewCommand creates a standardized cobra command.
func NewCommand[T any](cfg CommandConfig[T]) *cobra.Command {
	cmd := &cobra.Command{
		Use:               cfg.Use,
		Short:             cfg.Short,
		Long:              cfg.Long,
		Args:              cfg.Args,
		ValidArgsFunction: cfg.ValidArgsFunction,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.CtxFunc == nil {
				return nil
			}
			ctx := cfg.CtxFunc(cmd.Context(), cfg.Options)
			return cfg.Handler.Validate(ctx)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.CtxFunc == nil {
				ctx := cfg.Options
				if err := cfg.Handler.Execute(ctx); err != nil {
					return err
				}
				return cfg.Handler.Finalize(ctx)
			}
			ctx := cfg.CtxFunc(cmd.Context(), cfg.Options)
			if err := cfg.Handler.Execute(ctx); err != nil {
				return err
			}
			return cfg.Handler.Finalize(ctx)
		},
	}
	return cmd
}
