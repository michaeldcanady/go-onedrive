package shared

// LoginResult represents the outcome of a successful login operation.
type LoginResult struct {
	// Token is the access token acquired during login.
	Token AccessToken
	// ProfileName is the name of the configuration profile used for login.
	ProfileName string
}
