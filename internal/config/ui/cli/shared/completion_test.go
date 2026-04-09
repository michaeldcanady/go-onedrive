package shared

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestConfigKeyCompletion(t *testing.T) {
	// We don't need a real container as it's not used in ConfigKeyCompletion
	fn := ConfigKeyCompletion(nil)

	tests := []struct {
		name              string
		toComplete        string
		expectedResults   []string
		expectedDirective cobra.ShellCompDirective
	}{
		{
			name:              "empty toComplete returns top level prefixes",
			toComplete:        "",
			expectedResults:   []string{"auth.", "logging."},
			expectedDirective: cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace,
		},
		{
			name:              "aut returns auth.",
			toComplete:        "aut",
			expectedResults:   []string{"auth."},
			expectedDirective: cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace,
		},
		{
			name:              "auth. returns sub keys",
			toComplete:        "auth.",
			expectedResults:   []string{"auth.client_id", "auth.client_secret", "auth.method", "auth.provider", "auth.redirect_uri", "auth.tenant_id"},
			expectedDirective: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name:              "logging. returns sub keys",
			toComplete:        "logging.",
			expectedResults:   []string{"logging.format", "logging.level", "logging.output"},
			expectedDirective: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name:              "unknown prefix returns empty",
			toComplete:        "unknown",
			expectedResults:   nil,
			expectedDirective: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, directive := fn(nil, nil, tt.toComplete)
			assert.Equal(t, tt.expectedResults, results)
			assert.Equal(t, tt.expectedDirective, directive)
		})
	}
}
