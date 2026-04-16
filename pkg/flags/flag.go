package flags

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	flagTag = "flag"
)

// Flag represents the agnostic metadata for a CLI flag.
type Flag struct {
	Name        string
	Short       string
	Description string
	Default     string
	Persistent  bool
	Value       any          // Pointer to the struct field
	Kind        reflect.Kind // The underlying type of the field
}

// Parse extracts flag metadata from a struct using reflection.
func Parse(options any) ([]Flag, error) {
	v := reflect.ValueOf(options)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("options must be a struct or a pointer to a struct")
	}

	var parsedFlags []Flag
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(flagTag)
		if tagValue == "" {
			continue
		}

		f, err := parseField(v.Field(i), field, tagValue)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}
		parsedFlags = append(parsedFlags, f)
	}

	return parsedFlags, nil
}

func parseField(fieldValue reflect.Value, field reflect.StructField, tag string) (Flag, error) {
	// Simple state machine to parse: name,key1="value1",key2='value2'
	parts := parseTagParts(tag)
	if len(parts) == 0 {
		return Flag{}, fmt.Errorf("empty flag tag")
	}

	name := parts[0]
	f := Flag{
		Name:  name,
		Value: fieldValue.Addr().Interface(),
		Kind:  fieldValue.Kind(),
	}

	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := unquote(strings.TrimSpace(kv[1]))

		switch key {
		case "short":
			f.Short = value
		case "default":
			f.Default = value
		case "desc":
			f.Description = value
		case "persistent":
			f.Persistent = (value == "true")
		}
	}

	return f, nil
}

// parseTagParts splits the tag by comma, but respects quotes.
func parseTagParts(tag string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	var quoteChar rune

	runes := []rune(tag)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case (r == '"' || r == '\'') && !inQuote:
			inQuote = true
			quoteChar = r
			current.WriteRune(r)
		case r == quoteChar && inQuote:
			inQuote = false
			current.WriteRune(r)
		case r == ',' && !inQuote:
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

func unquote(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

// RegisterFlags is a convenience function that parses flags and registers them to Cobra.
func RegisterFlags(cmd *cobra.Command, options any) error {
	parsed, err := Parse(options)
	if err != nil {
		return err
	}

	return RegisterWithCobra(cmd, parsed)
}

// RegisterWithCobra is the Cobra-specific adapter for the agnostic Flag metadata.
func RegisterWithCobra(cmd *cobra.Command, parsedFlags []Flag) error {
	for _, f := range parsedFlags {
		flags := cmd.Flags()
		if f.Persistent {
			flags = cmd.PersistentFlags()
		}

		switch f.Kind {
		case reflect.Bool:
			def, _ := strconv.ParseBool(f.Default)
			flags.BoolVarP(f.Value.(*bool), f.Name, f.Short, def, f.Description)
		case reflect.String:
			flags.StringVarP(f.Value.(*string), f.Name, f.Short, f.Default, f.Description)
		case reflect.Int:
			def, _ := strconv.Atoi(f.Default)
			flags.IntVarP(f.Value.(*int), f.Name, f.Short, def, f.Description)
		case reflect.Slice:
			ptr, ok := f.Value.(*[]string)
			if !ok {
				return fmt.Errorf("unsupported slice type for flag %s", f.Name)
			}
			var def []string
			if f.Default != "" {
				def = strings.Split(f.Default, ";")
			}
			flags.StringSliceVarP(ptr, f.Name, f.Short, def, f.Description)
		default:
			return fmt.Errorf("unsupported type %s for flag %s", f.Kind, f.Name)
		}
	}
	return nil
}
