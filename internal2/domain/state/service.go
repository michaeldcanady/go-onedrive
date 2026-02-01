package state

type Service interface {
	GetCurrentProfile() (string, error)
	SetCurrentProfile(name string) error
	ClearCurrentProfile() error
}
