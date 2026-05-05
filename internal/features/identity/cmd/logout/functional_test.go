package logout

import (
	"bytes"
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout_Functional(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mCfg *mockConfigService, mIdent *mockIdentityService, mLog *mockLogger)
		wantErr bool
	}{
		{
			name: "successful logout flow",
			setup: func(mCfg *mockConfigService, mIdent *mockIdentityService, mLog *mockLogger) {
				mCfg.On("GetConfig", mock.Anything).Return(config.Config{
					Auth: config.AuthenticationConfig{Provider: "microsoft"},
				}, nil)
				mIdent.On("Logout", mock.Anything, "microsoft", "").Return(nil)
				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("Debug", mock.Anything, mock.Anything).Return()
				mLog.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mCfg := new(mockConfigService)
			mIdent := new(mockIdentityService)
			mLog := new(mockLogger)
			tt.setup(mCfg, mIdent, mLog)

			handler := NewCommand(mCfg, mIdent, mLog)

			buf := new(bytes.Buffer)
			cmdCtx := &CommandContext{
				Ctx: context.Background(),
				Options: &Options{
					Stdout: buf,
				},
			}

			err := handler.Validate(cmdCtx)
			assert.NoError(t, err)

			err = handler.Execute(cmdCtx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, buf.String(), "Logged out successfully.")
			}
		})
	}
}
