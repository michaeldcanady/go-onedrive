package get

import (
	"fmt"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

// Options provides the settings for the config get command.
type Options struct {
	// Key is the configuration key to retrieve (optional).
	Key string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Key != "" {
		keys := config.GetAvailableKeys()
		found := false
		for _, k := range keys {
			if strings.EqualFold(k, o.Key) {
				found = true
				break
			}
		}
		if !found {
			return errors.NewAppError(
				errors.CodeInvalidInput,
				fmt.Errorf("unknown configuration key: %s", o.Key),
				"unknown configuration key",
				"Use 'odc config get' without a key to see all available keys.",
			)
		}
	}
	return nil
}
