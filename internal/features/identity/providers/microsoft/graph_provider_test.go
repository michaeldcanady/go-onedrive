package microsoft

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTokenCredential struct {
	mock.Mock
}

func (m *mockTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(azcore.AccessToken), args.Error(1)
}

func TestGraphProvider_Name(t *testing.T) {
	p := NewGraphProvider(nil, nil)
	assert.Equal(t, "microsoft", p.Name())
}

func TestGraphProvider_Client(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mCred *mockTokenCredential, mLog *mockLogger) (azcore.TokenCredential, *mockLogger)
		wantErr bool
		errMsg  string
	}{
		{
			name: "initialization success",
			setup: func(mCred *mockTokenCredential, mLog *mockLogger) (azcore.TokenCredential, *mockLogger) {
				return mCred, mLog
			},
			wantErr: false,
		},
		{
			name: "missing credential",
			setup: func(mCred *mockTokenCredential, mLog *mockLogger) (azcore.TokenCredential, *mockLogger) {
				return nil, mLog
			},
			wantErr: true,
			errMsg:  "no authentication credential provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mCred := new(mockTokenCredential)
			mLog := new(mockLogger)

			cred, log := tt.setup(mCred, mLog)
			p := NewGraphProvider(cred, log)

			client, err := p.Client(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				// Second call should return the same client
				client2, err := p.Client(ctx)
				assert.NoError(t, err)
				assert.Same(t, client, client2)
			}
		})
	}
}
