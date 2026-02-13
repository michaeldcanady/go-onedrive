package msalclient

import (
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func authorityFromTenant(tenantID string) string {
	if tenantID == "" || tenantID == "common" {
		return "https://login.microsoftonline.com/common"
	}
	return fmt.Sprintf("https://login.microsoftonline.com/%s", tenantID)
}

func (f *Factory) PublicClient(cfg config.AuthenticationConfigImpl) (public.Client, error) {
	return public.New(cfg.ClientID,
		public.WithAuthority(authorityFromTenant(cfg.TenantID)),
	)
}

func (f *Factory) ConfidentialClient(cfg config.AuthenticationConfigImpl) (confidential.Client, error) {
	cred, err := confidential.NewCredFromSecret(cfg.ClientSecret)
	if err != nil {
		return confidential.Client{}, err
	}

	return confidential.New(authorityFromTenant(cfg.TenantID), cfg.ClientID, cred)
}
