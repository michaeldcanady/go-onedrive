package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type TestOptions struct {
	Verbose    bool     `flag:"verbose,short=v,desc='Enable verbose logging',default=false"`
	Count      int      `flag:"count,short=c,desc=\"Number of items\",default=10"`
	Name       string   `flag:"name,short=n,desc='User name',default='guest'"`
	Tags       []string `flag:"tags,short=t,desc='List of tags',default='a;b;c'"`
	Persistent bool     `flag:"global,desc='Global flag',persistent=true,default=true"`
	WithComma  string   `flag:"comma,desc='Description, with comma',default='default,value'"`
}

func TestRegisterFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	opts := &TestOptions{}

	err := RegisterFlags(cmd, opts)
	assert.NoError(t, err)

	// Check flags
	v, err := cmd.Flags().GetBool("verbose")
	assert.NoError(t, err)
	assert.False(t, v)

	c, err := cmd.Flags().GetInt("count")
	assert.NoError(t, err)
	assert.Equal(t, 10, c)

	n, err := cmd.Flags().GetString("name")
	assert.NoError(t, err)
	assert.Equal(t, "guest", n)

	tags, err := cmd.Flags().GetStringSlice("tags")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, tags)

	g, err := cmd.PersistentFlags().GetBool("global")
	assert.NoError(t, err)
	assert.True(t, g)

	commaDesc := cmd.Flags().Lookup("comma").Usage
	assert.Equal(t, "Description, with comma", commaDesc)

	commaVal, err := cmd.Flags().GetString("comma")
	assert.NoError(t, err)
	assert.Equal(t, "default,value", commaVal)

	// Test flag changes
	err = cmd.Flags().Set("verbose", "true")
	assert.NoError(t, err)
	assert.True(t, opts.Verbose)

	err = cmd.Flags().Set("count", "42")
	assert.NoError(t, err)
	assert.Equal(t, 42, opts.Count)
}
