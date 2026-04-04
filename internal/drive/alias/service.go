package alias

type Service interface {
	// GetDriveIDByAlias retrieves the drive ID associated with a given alias, if it exists.
	GetDriveIDByAlias(alias string) (string, error)
	// GetAliasByDriveID retrieves the alias for a given drive ID, if it exists.
	GetAliasByDriveID(driveID string) (string, error)
	// SetAlias assigns an alias to a specific drive ID.
	SetAlias(driveID string, alias string) error
	// DeleteAlias removes the alias.
	DeleteAlias(alias string) error
	// ListAliases returns a mapping of all drive IDs to their respective aliases.
	ListAliases() (map[string]string, error)
}
