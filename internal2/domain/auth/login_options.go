package auth

type LoginOptions struct {
	Force     bool
	Scopes    []string
	EnableCAE bool
}
