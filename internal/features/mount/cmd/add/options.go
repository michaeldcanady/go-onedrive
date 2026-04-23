package add

import "io"

type Options struct {
	MountOptions []string
	IdentityID   string
	StdErr       io.Writer
	Stdout       io.Writer
}

func NewOptions() *Options {
	return &Options{
		MountOptions: make([]string, 0),
	}
}
