package app

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type AuthService interface {
	Login(ctx context.Context) error
}

type AuthServiceImpl struct {
	config config.Config
	client *msgraphsdkgo.GraphServiceClient
}

func NewAuthService(config config.Config, client *msgraphsdkgo.GraphServiceClient) *AuthServiceImpl {
	return &AuthServiceImpl{
		config: config,
		client: client,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context) error {
	return nil
}
