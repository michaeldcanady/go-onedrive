package state

type State struct {
	CurrentProfile string            `yaml:"currentProfile" json:"current_profile"`
	CurrentDrive   string            `yaml:"currentDrive" json:"current_drive"`
	DriveAliases   map[string]string `yaml:"driveAliases" json:"drive_aliases"`
}
