package auth

type Status struct {
	LoggedIn     bool
	Username     string
	TenantID     string
	ExpiresInSec int64
}
