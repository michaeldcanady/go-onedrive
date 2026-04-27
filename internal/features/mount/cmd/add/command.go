package add

import (
	"fmt"
	"maps"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/spf13/cobra"
)

// baseFlagCompletion builds option completion
func baseFlagCompletion(mountOptions map[string][]mount.MountOption) cobra.CompletionFunc {
	backendOpts := buildBackendOpts(mountOptions)

	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) < 2 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		backendType := args[1]
		opts, ok := backendOpts[backendType]
		if !ok {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if strings.Contains(toComplete, "=") {
			return completeOptionValue(toComplete, opts)
		}

		return completeOptionKey(toComplete, opts), cobra.ShellCompDirectiveNoFileComp
	}
}

type optInfo struct {
	values    []string
	validator func(string) bool
}

func buildBackendOpts(mountOptions map[string][]mount.MountOption) map[string]map[string]optInfo {
	backendOpts := map[string]map[string]optInfo{}
	for backend, opts := range mountOptions {
		backendOpts[backend] = make(map[string]optInfo)
		for _, opt := range opts {
			o := opt
			backendOpts[backend][strings.ToLower(opt.Key)] = optInfo{
				values:    o.Values,
				validator: createValidator(o.Values),
			}
		}
	}
	return backendOpts
}

func createValidator(values []string) func(string) bool {
	return func(value string) bool {
		if len(values) == 0 {
			return true
		}
		for _, v := range values {
			if v == value {
				return true
			}
		}
		return false
	}
}

func completeOptionValue(toComplete string, opts map[string]optInfo) ([]cobra.Completion, cobra.ShellCompDirective) {
	components := strings.SplitN(toComplete, "=", 2)
	key := strings.ToLower(components[0])
	valuePrefix := components[1]

	info, ok := opts[key]
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completes []string
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

func completeOptionKey(toComplete string, opts map[string]optInfo) []string {
	var completes []string
	for key := range maps.Keys(opts) {
		if strings.HasPrefix(key, strings.ToLower(toComplete)) {
			completes = append(completes, key+"=")
		}
	}
	return completes
}

func CreateAddCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("mount-add")
	handler := NewCommand(container.Mounts(), container.Identity(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:   "add <path> <type> <identity_id>",
		Short: "Add a mount point",
		Args:  cobra.ExactArgs(3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Type = args[1]
			opts.IdentityID = args[2]
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c = NewCommandContext(cmd.Context(), &opts)

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	if err := cmd.RegisterFlagCompletionFunc("option", baseFlagCompletion(container.Mounts().GetMountOptions())); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "failed to register flag completion: %v\n", err)
	}

	cmd.Flags().StringSliceVar(&opts.MountOptions, "option", []string{}, "Provider-specific options in key=value format (repeatable)")

	return cmd
}
