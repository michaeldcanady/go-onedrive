package filtering

import (
	"errors"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type FilterOptions struct {
	ItemType   domainfs.ItemType
	IncludeAll bool
}

func NewFilterOptions() *FilterOptions {
	return &FilterOptions{
		ItemType: domainfs.ItemTypeUnknown,
	}
}

func (o *FilterOptions) Apply(opts []FilterOption) error {
	if o == nil {
		return errors.New("config is nil")
	}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}
	return nil
}

type FilterOption = func(*FilterOptions) error

func WithItemType(itemType domainfs.ItemType) FilterOption {
	if itemType == domainfs.ItemTypeUnknown {
		return func(_ *FilterOptions) error {
			return errors.New("filtered item type is unknown")
		}
	}
	return func(config *FilterOptions) error {
		config.ItemType = itemType
		return nil
	}
}

func IncludeAll() FilterOption {
	return func(config *FilterOptions) error {
		config.IncludeAll = true
		return nil
	}
}

func ExcludeHidden() FilterOption {
	return func(config *FilterOptions) error {
		config.IncludeAll = false
		return nil
	}
}
