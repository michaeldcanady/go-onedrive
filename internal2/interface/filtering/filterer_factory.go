package filtering

import "errors"

type FilterFactory struct{}

func NewFilterFactory() *FilterFactory {
	return &FilterFactory{}
}

func (f *FilterFactory) Create(filterType string) (Filterer, error) {
	switch filterType {
	case "hidden":
		return NewHiddenFilterer(), nil
	case "none":
		return NewNoOpFilter(), nil
	}

	return nil, errors.New("invalid filterType")
}
