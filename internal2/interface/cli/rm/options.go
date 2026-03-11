package rm

import (
	"errors"

	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Options struct {
	util.BaseOptions

	Path      string
	Permanent bool
	Force     bool
}

func (o Options) Validate() error {
	if o.Path == "" {
		return errors.New("path is required")
	}
	return nil
}
