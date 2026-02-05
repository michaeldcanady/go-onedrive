package filtering

import "reflect"

type Service interface {
	Register(typ string, filter Filter) error
	RegisterWithType(format string, typ reflect.Type, f Filter) error
	Filter(typ string, v any) error
}
