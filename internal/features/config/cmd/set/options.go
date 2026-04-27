package set

import (
	"errors"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/pkg/validation"
)

// Options provides the settings for the config set command.
type Options struct {
	// Key is the configuration key to update.
	Key string
	// Value is the new value for the configuration setting.
	Value string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer

	Stderr io.Writer
}

func NewOptions() *Options {
	return &Options{}
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	p := validation.All(
		validation.Required(func(o Options) string { return o.Key }, "key"),
		validation.Required(func(o Options) string { return o.Value }, "value"),
		validation.PolicyFunc[Options](func(o Options) error {
			cleanKey := strings.TrimSpace(o.Key)
			illegalChars := []string{" ", "\n", "\r"}
			for _, char := range illegalChars {
				if strings.Contains(cleanKey, char) {
					return errors.New("key contains illegal characters")
				}
			}
			return nil
		}),
	)

	return p.Evaluate(*o)
}
