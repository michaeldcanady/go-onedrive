package formatting

import (
	"io"
	"reflect"
)

type Service interface {
	Format(w io.Writer, format string, v any) error
	Register(format string, formatter Formatter) error
	RegisterWithType(format string, typ reflect.Type, f Formatter) error
}
