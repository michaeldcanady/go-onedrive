package domain

type Options struct {
	Method      Method
	TenantID    string
	ClientID    string
	LoginHint   string
	RedirectURL string
}
