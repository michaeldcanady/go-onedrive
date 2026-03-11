package domain

type Service interface {
	GetCurrentProfile() (string, error)
	SetCurrentProfile(name string, scope Scope) error
	ClearCurrentProfile() error

	GetCurrentDrive() (string, error)
	SetCurrentDrive(id string, scope Scope) error
	ClearCurrentDrive() error

	GetDriveAlias(alias string) (string, error)
	SetDriveAlias(alias, driveID string) error
	RemoveDriveAlias(alias string) error
	ListDriveAliases() (map[string]string, error)
}
