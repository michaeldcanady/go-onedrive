package domain

type Service interface {
	GetCurrentProfile() (string, error)
	SetCurrentProfile(name string) error
	ClearCurrentProfile() error
	SetSessionProfile(name string)

	GetCurrentDrive() (string, error)
	SetCurrentDrive(name string) error
	ClearCurrentDrive() error
	SetSessionDrive(driveID string)

	GetDriveAlias(alias string) (string, error)
	SetDriveAlias(alias, driveID string) error
	RemoveDriveAlias(alias string) error
	ListDriveAliases() (map[string]string, error)
}
