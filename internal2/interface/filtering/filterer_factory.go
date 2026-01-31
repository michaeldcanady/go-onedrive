package filtering

type FilterFactory struct{}

func NewFilterFactory() *FilterFactory {
	return &FilterFactory{}
}

func (f *FilterFactory) Create(opts ...FilterOption) (Filterer, error) {
	config := NewFilterOptions()
	if err := config.Apply(opts); err != nil {
		return nil, err
	}

	return NewOptionsFilterer(*config), nil
}
