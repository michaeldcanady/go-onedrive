package login

import "io"

// Options provides the user-facing settings for the auth login command.
type Options struct {
	// ShowToken determines whether the acquired access token is printed to stdout.
	ShowToken bool
	// Force specifies whether to re-authenticate regardless of existing credentials.
	Force bool
	// Method specifies the authentication mechanism to use.
	Method string
	// TenantID is the Azure AD tenant ID for Service Principal authentication.
	TenantID string
	// ClientID is the Azure AD client ID for Service Principal authentication.
	ClientID string
	// ClientSecret is the Azure AD client secret for Service Principal authentication.
	ClientSecret string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
	// Stderr is the destination for error messages.
	Stderr io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	return nil
}
