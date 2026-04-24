package mount

// OptionsProvider defines the interface for backends that support options
type OptionsProvider interface {
	// ProvideOptions returns supported options for the backends
	ProvideOptions() []MountOption
}
