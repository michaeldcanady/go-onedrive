package account

type Account struct {
	// Authority is the URL of the authority that issued the token.
	Authority string `json:"authority"`

	// ClientID is the ID of the application that authenticated the user.
	ClientID string `json:"clientId"`

	// HomeAccountID uniquely identifies the account.
	HomeAccountID string `json:"homeAccountId"`

	// TenantID identifies the tenant in which the user authenticated.
	TenantID string `json:"tenantId"`

	// Username is the user's preferred username.
	Username string `json:"username"`

	// Version of the AuthenticationRecord.
	Version string `json:"version"`
}
