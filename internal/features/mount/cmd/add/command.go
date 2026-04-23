package add

import (
	"context"
	"maps"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/spf13/cobra"
)

type noOpAccountGetter struct{}

func (_ *noOpAccountGetter) GetAccount(ctx context.Context, identityID string) (*identity.Account, error) {
	return &identity.Account{
		ID: identityID,
	}, nil
}

// baseFlagCompletion builds option completion
// map is option = slice of options
func baseFlagCompletion(mountOptions map[string][]mount.MountOption) cobra.CompletionFunc {
	type optInfo struct {
		values    []string
		validator func(string) bool
	}
	backendOpts := map[string]map[string]optInfo{}

	for backend, opts := range mountOptions {
		backendOpts[backend] = make(map[string]optInfo)
		for _, opt := range opts {
			o := opt // capture loop var
			backendOpts[backend][strings.ToLower(opt.Key)] = optInfo{
				values: opt.Values,
				validator: func(value string) bool {
					if len(o.Values) == 0 {
						return true
					}
					for _, v := range o.Values {
						if v == value {
							return true
						}
					}
					return false
				},
			}
		}
	}

	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) < 2 {
			// Not enough args to determine backend type
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		backendType := args[1]
		opts, ok := backendOpts[backendType]
		if !ok {
			// no opts supported for this backend
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if strings.Contains(toComplete, "=") {
			components := strings.SplitN(toComplete, "=", 2)
			key := strings.ToLower(components[0])
			valuePrefix := components[1]

			info, ok := opts[key]
			if !ok {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completes []cobra.Completion
			for _, v := range info.values {
				if strings.HasPrefix(v, valuePrefix) {
					completes = append(completes, components[0]+"="+v)
				}
			}

			if len(completes) == 0 && info.validator != nil && info.validator(valuePrefix) {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return completes, cobra.ShellCompDirectiveNoFileComp
		}

		var completes []cobra.Completion
		for key := range maps.Keys(opts) {
			if strings.HasPrefix(key, strings.ToLower(toComplete)) {
				completes = append(completes, key+"=")
			}
		}
		return completes, cobra.ShellCompDirectiveNoFileComp
	}
}

func CreateAddCmd(container di.Container) *cobra.Command {
	opts := NewOptions()

	l, _ := container.Logger().CreateLogger("mount-add")
	// TODO: need really account getter
	handler := NewCommand(container.Mounts(), &noOpAccountGetter{}, l)

	cmd := cli.NewCommand(cli.CommandConfig[CommandContext]{
		Use:     "add <path> <type> <identity_id>",
		Short:   "Add a mount point",
		Args:    cobra.ExactArgs(3),
		Handler: handler,
		Options: NewCommandContext(context.Background(), opts),
		CtxFunc: func(ctx context.Context, o *CommandContext) *CommandContext {
			o.Ctx = ctx
			return o
		},
	})

	cmd.RegisterFlagCompletionFunc("option", baseFlagCompletion(container.Mounts().GetMountOptions()))

	cmd.Flags().StringSliceVar(&opts.MountOptions, "option", []string{}, "Provider-specific options in key=value format (repeatable)")

	return cmd
}
