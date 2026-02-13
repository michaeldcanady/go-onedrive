package config

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

var _ AuthenticationConfig = (*AuthenticationConfigImpl)(nil)

type AuthenticationConfig interface {
	GetAuthenticationType() string
	GetClientID() string
	GetTenantID() string
	GetRedirectURI() string
	GetAuthenticationRecord() *azidentity.AuthenticationRecord
	GetProfileCache() string
}

type AuthenticationConfigImpl struct {
	Type         string `mapstructure:"type" yaml:"type"`
	ClientID     string `mapstructure:"client_id" yaml:"client_id"`
	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`
	TenantID     string `mapstructure:"tenant_id" yaml:"tenant_id"`
	RedirectURI  string `mapstructure:"redirect_uri" yaml:"redirect_uri"`
	ProfileCache string `mapstructure:"profile_cache" yaml:"profile_cache"`
	// TODO: move some where else
	AuthenticationRecord *azidentity.AuthenticationRecord
}

// GetAuthenticationRecord implements AuthenticationConfig.
func (a *AuthenticationConfigImpl) GetAuthenticationRecord() *azidentity.AuthenticationRecord {
	return a.AuthenticationRecord
}

func (a *AuthenticationConfigImpl) GetAuthenticationType() string {
	return a.Type
}

func (a *AuthenticationConfigImpl) GetClientID() string {
	return a.ClientID
}

func (a *AuthenticationConfigImpl) GetTenantID() string {
	return a.TenantID
}

func (a *AuthenticationConfigImpl) GetRedirectURI() string {
	return a.RedirectURI
}

func (a *AuthenticationConfigImpl) GetProfileCache() string {
	return a.ProfileCache
}
