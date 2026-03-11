package domain

type LoginOptions struct {
	Force     bool
	Scopes    []string
	EnableCAE bool
}
