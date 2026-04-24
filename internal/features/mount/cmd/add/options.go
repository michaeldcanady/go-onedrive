package add

import "io"

type Options struct {
	MountOptions []string
	IdentityID   string
	Path         string
	Type         string
	Stderr       io.Writer
	Stdout       io.Writer
}
