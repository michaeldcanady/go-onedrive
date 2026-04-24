package mount

// OptionValidator defines the interface for backends that support option validation.
type OptionValidator interface {
	// ValidateOptions validates the provided options
	ValidateOptions(opts map[string]string) error
}
