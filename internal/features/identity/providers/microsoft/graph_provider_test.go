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
	ctx := context.Background()
	mCred := new(mockTokenCredential)
	mLog := new(mockLogger)

	p := NewGraphProvider(mCred, mLog)

	t.Run("initialization success", func(t *testing.T) {
		client, err := p.Client(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		
		// Second call should return the same client
		client2, err := p.Client(ctx)
		assert.NoError(t, err)
		assert.Same(t, client, client2)
	})

	t.Run("missing credential", func(t *testing.T) {
		pNoCred := NewGraphProvider(nil, mLog)
		_, err := pNoCred.Client(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no authentication credential provided")
	})
}
