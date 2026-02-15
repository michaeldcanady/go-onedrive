package auth

type Options struct {
	Method      Method
	TenantID    string
	ClientID    string
	LoginHint   string
	RedirectURL string
}
