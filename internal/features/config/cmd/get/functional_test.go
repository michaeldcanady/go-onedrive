package get

import (
	"bytes"
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigGet_Functional(t *testing.T) {
	ctx := context.Background()

	mSvc := new(mockConfigService)
	mLog := new(mockLogger)

	mSvc.On("GetConfig", mock.Anything).Return(config.Config{
		Auth: config.AuthenticationConfig{Provider: "functional-test"},
	}, nil)

	mLog.On("WithContext", mock.Anything).Return(mLog)
	mLog.On("Debug", mock.Anything, mock.Anything).Return()

	handler := NewCommand(mSvc, mLog)

	buf := new(bytes.Buffer)
	cmdCtx := &CommandContext{
		Ctx: ctx,
		Options: &Options{
			Key:    "auth.provider",
			Stdout: buf,
		},
	}

	err := handler.Validate(cmdCtx)
	assert.NoError(t, err)

	err = handler.Execute(cmdCtx)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "auth.provider: functional-test")
}
