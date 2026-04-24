package config

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Config represents the application's top-level configuration structure.
type Config struct {
	// Auth contains settings related to authentication and provider identity.
	Auth AuthenticationConfig `json:"auth" yaml:"auth"`
	// Logging contains settings related to logging behavior and output.
	Logging LoggingConfig `json:"logging" yaml:"logging"`
	// Mounts defines the collection of virtual filesystem mount points.
	Mounts []MountConfig `json:"mounts,omitempty" yaml:"mounts,omitempty" mapstructure:"mounts"`
	// Editor contains settings related to the external editor service.
	Editor EditorConfig `json:"editor" yaml:"editor"`
}

// SetValue updates a field in the config based on the provided key (e.g., "auth.provider").
func (c *Config) SetValue(key, value string) error {
	return setFieldRecursive(reflect.ValueOf(c).Elem(), strings.Split(key, "."), value)
}

func setFieldRecursive(v reflect.Value, path []string, val string) error {
	if len(path) == 0 {
		return fmt.Errorf("no field specified")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == path[0] {
			target := v.Field(i)
			if len(path) == 1 {
				return setReflectValue(target, val)
			}
			return setFieldRecursive(target, path[1:], val)
		}
	}
	return fmt.Errorf("field %s not found", path[0])
}

func setReflectValue(target reflect.Value, val string) error {
	// Check if the type implements encoding.TextUnmarshaler
	if target.CanInterface() {
		if unmarshaler, ok := target.Interface().(encoding.TextUnmarshaler); ok {
			return unmarshaler.UnmarshalText([]byte(val))
		}
		// Also check if a pointer to the type implements it
		if target.CanAddr() {
			if unmarshaler, ok := target.Addr().Interface().(encoding.TextUnmarshaler); ok {
				return unmarshaler.UnmarshalText([]byte(val))
			}
		}
	}

	switch target.Kind() {
	case reflect.String:
		target.SetString(val)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", val)
		}
		target.SetBool(b)
	case reflect.Int, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", val)
		}
		target.SetInt(i)
	default:
		return fmt.Errorf("unsupported type: %s", target.Kind())
	}
	return nil
}
