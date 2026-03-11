package domain

type Service interface {
	Get(key Key) (string, error)
	Set(key Key, value string, scope Scope) error
	Clear(key Key) error

	GetDriveAlias(alias string) (string, error)
	SetDriveAlias(alias, driveID string) error
	RemoveDriveAlias(alias string) error
	ListDriveAliases() (map[string]string, error)
}
