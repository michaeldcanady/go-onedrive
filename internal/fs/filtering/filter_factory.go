package filtering

// FilterFactory provides operations for initializing configured filterer instances.
type FilterFactory struct{}

// NewFilterFactory initializes a new instance of the FilterFactory.
func NewFilterFactory() *FilterFactory {
	return &FilterFactory{}
}

// Create initializes and configures a new filterer with the provided functional options.
func (f *FilterFactory) Create(opts ...FilterOption) (Filterer, error) {
	config := NewFilterOptions()
	if err := config.Apply(opts); err != nil {
		return nil, err
	}

	return NewOptionsFilterer(*config), nil
}
