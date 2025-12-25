package config

type AuthenticationConfig interface {
	GetAuthenticationType() string
	GetClientID() string
	GetTenantID() string
	GetRedirectURI() string
}

type AuthenticationConfigImpl struct {
	Type        string `mapstructure:"type"`
	ClientID    string `mapstructure:"client_id"`
	TenantID    string `mapstructure:"tenant_id"`
	RedirectURI string `mapstructure:"redirect_uri"`
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
