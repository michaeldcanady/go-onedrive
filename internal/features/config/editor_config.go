package config

// EditorConfig represents the configuration for the external editor.
type EditorConfig struct {
	// Command is the explicit editor command to use.
	Command string `json:"command" yaml:"command"`
}
