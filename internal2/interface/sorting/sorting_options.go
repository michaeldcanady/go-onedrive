package sorting

import "errors"

type SortingOptions struct {
	Direction Direction
	Field     string
}

func NewSortingOptions() *SortingOptions {
	return &SortingOptions{
		Direction: DirectionAscending,
		Field:     "",
	}
}

func (o *SortingOptions) Apply(opts ...SortingOption) error {
	if o == nil {
		return errors.New("nil options")
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}

type SortingOption = func(*SortingOptions) error

func WithDirection(direction Direction) SortingOption {
	return func(opt *SortingOptions) error {
		if opt == nil {
			return errors.New("opt is nil")
		}
		opt.Direction = direction
		return nil
	}
}

func WithField(name string) SortingOption {
	return func(opt *SortingOptions) error {
		if opt == nil {
			return errors.New("opt is nil")
		}
		opt.Field = name
		return nil
	}
}
