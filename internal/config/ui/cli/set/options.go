package set

import (
	"fmt"
	"io"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

// Options provides the settings for the config set command.
type Options struct {
	// Key is the configuration key to set.
	Key string
	// Value is the configuration value to set.
	Value string
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	if o.Key == "" {
		return errors.NewAppError(
			errors.CodeInvalidInput,
			nil,
			"configuration key is required",
			"Please provide a key to set.",
		)
	}

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
			"Use 'odc config get' to see all available keys.",
		)
	}

	// Validate allowed values if the schema provides them (enums/anyOf consts)
	allowed := config.GetAllowedValues(o.Key)
	if len(allowed) > 0 {
		validValue := false
		for _, v := range allowed {
			if strings.EqualFold(v, o.Value) {
				validValue = true
				break
			}
		}
		if !validValue && config.IsStrictEnum(o.Key) {
			return errors.NewAppError(
				errors.CodeInvalidInput,
				fmt.Errorf("invalid value for key %s: %s", o.Key, o.Value),
				"invalid configuration value",
				fmt.Sprintf("Allowed values for '%s' are: %s", o.Key, strings.Join(allowed, ", ")),
			)
		}
	}

	return nil
}
