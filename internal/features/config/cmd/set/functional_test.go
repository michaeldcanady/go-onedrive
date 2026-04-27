package set

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigSet_Functional(t *testing.T) {
	ctx := context.Background()

	mSvc := new(mockConfigService)
	mLog := new(mockLogger)

	mSvc.On("UpdateConfig", mock.Anything, "logging.level", "info").Return(nil)

	mLog.On("WithContext", mock.Anything).Return(mLog)
	mLog.On("Debug", mock.Anything, mock.Anything).Return()
	mLog.On("Info", mock.Anything, mock.Anything).Return()

	handler := NewCommand(mSvc, mLog)

	buf := new(bytes.Buffer)
	cmdCtx := &CommandContext{
		Ctx: ctx,
		Options: &Options{
			Key:    "logging.level",
			Value:  "info",
			Stdout: buf,
		},
	}

	err := handler.Validate(cmdCtx)
	assert.NoError(t, err)

	err = handler.Execute(cmdCtx)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Set logging.level to info successfully")
}
