package args

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestOptions struct {
	Source      string `arg:"1"`
	Destination string `arg:"2"`
}

func TestBind(t *testing.T) {
	opts := &TestOptions{}
	inputArgs := []string{"foo", "bar"}

	err := Bind(inputArgs, opts)
	assert.NoError(t, err)
	assert.Equal(t, "foo", opts.Source)
	assert.Equal(t, "bar", opts.Destination)
}

func TestExactArgs(t *testing.T) {
	opts := &TestOptions{}
	validator := ExactArgs(opts)
	// We need to pass a nil cmd to check the validator
	assert.NoError(t, validator(nil, []string{"a", "b"}))
	assert.Error(t, validator(nil, []string{"a"}))
}
