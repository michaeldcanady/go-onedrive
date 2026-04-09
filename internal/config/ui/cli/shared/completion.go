package shared

import (
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// ConfigKeyCompletion returns a completion function for configuration keys.
func ConfigKeyCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		keys := config.GetAvailableKeys()

		prefixes := make(map[string]struct{})
		var results []string

		depth := strings.Count(toComplete, ".")

		for _, key := range keys {
			if startsWithIgnoreCase(key, toComplete) {
				parts := strings.Split(key, ".")
				if depth < len(parts)-1 {
					// Add prefix for next level
					prefix := strings.Join(parts[:depth+1], ".") + "."
					if _, ok := prefixes[prefix]; !ok {
						results = append(results, prefix)
						prefixes[prefix] = struct{}{}
					}
				} else {
					// Add full key
					results = append(results, key)
				}
			}
		}

		directive := cobra.ShellCompDirectiveNoFileComp
		for _, res := range results {
			if strings.HasSuffix(res, ".") {
				directive |= cobra.ShellCompDirectiveNoSpace
				break
			}
		}

		sort.Strings(results)
		return results, directive
	}
}

// ConfigValueCompletion returns a completion function for configuration values based on the provided key.
func ConfigValueCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		key := args[0]
		values := config.GetAllowedValues(key)
		var results []string
		for _, val := range values {
			if toComplete == "" || startsWithIgnoreCase(val, toComplete) {
				results = append(results, val)
			}
		}

		directive := cobra.ShellCompDirectiveNoFileComp
		if key == "logging.output" {
			directive = cobra.ShellCompDirectiveDefault
		}

		return results, directive
	}
}

// SetCompletion returns a completion function for the 'config set' command.
func SetCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return ConfigKeyCompletion(container)(cmd, args, toComplete)
		}
		if len(args) == 1 {
			return ConfigValueCompletion(container)(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func startsWithIgnoreCase(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}
