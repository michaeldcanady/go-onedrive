package state

type Service interface {
	GetCurrentProfile() (string, error)
	SetCurrentProfile(name string) error
	ClearCurrentProfile() error
	SetSessionProfile(name string)

	GetCurrentDrive() (string, error)
	SetCurrentDrive(name string) error
	ClearCurrentDrive() error
}
