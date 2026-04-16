package args

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
)

const argTag = "arg"

// Bind populates the target struct fields with values from the args slice
// based on the "arg" struct tag (1-based index).
func Bind(args []string, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct or a pointer to a struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(argTag)
		if tag == "" {
			continue
		}

		index, err := strconv.Atoi(tag)
		if err != nil {
			return fmt.Errorf("invalid arg tag on field %s: %w", field.Name, err)
		}

		// Adjust 1-based tag to 0-based slice index
		idx := index - 1
		if idx < 0 {
			return fmt.Errorf("arg tag must be 1 or greater on field %s", field.Name)
		}

		if idx < len(args) {
			fieldValue := v.Field(i)
			if !fieldValue.CanSet() {
				continue
			}

			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(args[idx])
			} else {
				return fmt.Errorf("unsupported arg field type: %s (only string supported currently)", fieldValue.Kind())
			}
		}
	}

	return nil
}

// ExactArgs returns a cobra.PositionalArgs validator based on the number of "arg" tags found.
func ExactArgs(target any) cobra.PositionalArgs {
	maxArg := countArgs(target)
	return cobra.ExactArgs(maxArg)
}

// MaximumNArgs returns a cobra.PositionalArgs validator based on the number of "arg" tags found.
func MaximumNArgs(target any) cobra.PositionalArgs {
	maxArg := countArgs(target)
	return cobra.MaximumNArgs(maxArg)
}

func countArgs(target any) int {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	maxArg := 0
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(argTag)
		if tag == "" {
			continue
		}
		index, _ := strconv.Atoi(tag)
		if index > maxArg {
			maxArg = index
		}
	}
	return maxArg
}
