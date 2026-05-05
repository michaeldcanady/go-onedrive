package microsoft

import (
	"context"
	"testing"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountStore struct {
	mock.Mock
}

func (m *mockAccountStore) Save(ctx context.Context, provider string, token identity.AccessToken) error {
	return m.Called(ctx, provider, token).Error(0)
}

func (m *mockAccountStore) Get(ctx context.Context, provider string, accountID string) (identity.AccessToken, error) {
	args := m.Called(ctx, provider, accountID)
	return args.Get(0).(identity.AccessToken), args.Error(1)
}

func (m *mockAccountStore) List(ctx context.Context, provider string) ([]string, error) {
	args := m.Called(ctx, provider)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockAccountStore) Delete(ctx context.Context, provider string, accountID string) error {
	return m.Called(ctx, provider, accountID).Error(0)
}

func (m *mockAccountStore) Close() error {
	return m.Called().Error(0)
}

func TestMicrosoftAuthorizer_Token(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name      string
		accountID string
		setup     func(m *mockAccountStore)
		wantErr   bool
		wantToken string
	}{
		{
			name:      "success",
			accountID: "user1",
			setup: func(m *mockAccountStore) {
				m.On("Get", mock.Anything, "microsoft", "user1").Return(identity.AccessToken{
					Token:     "test-token",
					ExpiresAt: now,
					Scopes:    []string{"scope1"},
				}, nil)
			},
			wantErr:   false,
			wantToken: "test-token",
		},
		{
			name:      "store error",
			accountID: "user1",
			setup: func(m *mockAccountStore) {
				m.On("Get", mock.Anything, "microsoft", "user1").Return(identity.AccessToken{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mStore := new(mockAccountStore)
			tt.setup(mStore)

			auth := NewMicrosoftAuthorizer(mStore)
			resp, err := auth.Token(ctx, &proto.GetTokenRequest{IdentityId: tt.accountID})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, resp.GetToken().GetToken())
			}
			mStore.AssertExpectations(t)
		})
	}
}
