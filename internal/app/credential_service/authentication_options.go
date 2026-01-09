package credentialservice

type authenticationConfig struct {
	// Claims are any additional claims required for the token to satisfy a conditional access policy, such as a
	// service may return in a claims challenge following an authorization failure. If a service returned the
	// claims value base64 encoded, it must be decoded before setting this field.
	Claims string

	// EnableCAE indicates whether to enable Continuous Access Evaluation (CAE) for the requested token. When true,
	// azidentity credentials request CAE tokens for resource APIs supporting CAE. Clients are responsible for
	// handling CAE challenges. If a client that doesn't handle CAE challenges receives a CAE token, it may end up
	// in a loop retrying an API call with a token that has been revoked due to CAE.
	EnableCAE bool

	// Scopes contains the list of permission scopes required for the token.
	Scopes []string

	// TenantID identifies the tenant from which to request the token. azidentity credentials authenticate in
	// their configured default tenants when this field isn't set.
	TenantID string

	Force bool
}

type AuthenticationOption = func(*authenticationConfig) error

func WithTenantID(tenantID string) AuthenticationOption {
	return func(cfg *authenticationConfig) error {
		cfg.TenantID = tenantID
		return nil
	}
}

func WithScopes(scopes ...string) AuthenticationOption {
	return func(cfg *authenticationConfig) error {
		cfg.Scopes = scopes
		return nil
	}
}

func WithClaims(claims string) AuthenticationOption {
	return func(cfg *authenticationConfig) error {
		cfg.Claims = claims
		return nil
	}
}

func WithCAE() AuthenticationOption {
	return func(cfg *authenticationConfig) error {
		cfg.EnableCAE = true
		return nil
	}
}

func WithForceAuthentication() AuthenticationOption {
	return func(cfg *authenticationConfig) error {
		cfg.Force = true
		return nil
	}
}

func buildConfig[T any](config T, opts ...func(T) error) error {
	// TODO: check if config is a pointer

	for _, opt := range opts {
		if err := opt(config); err != nil {
			return err
		}
	}
	return nil
}
