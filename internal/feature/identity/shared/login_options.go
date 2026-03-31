package shared

// LoginOptions provides settings for an authentication request.
type LoginOptions struct {
	// Force specifies whether to re-authenticate regardless of existing credentials.
	Force bool
	// Interactive specifies whether a UI interaction (like a browser) is allowed.
	Interactive bool
	// Method specifies the mechanism used for authentication.
	Method AuthMethod
	// ProviderSpecific contains extra parameters for the provider.
	ProviderSpecific map[string]string
}
