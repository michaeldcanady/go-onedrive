package logger

import (
	"fmt"
	"time"
)

// String creates a Field with a string value.
func String(key, val string) Field {
	return Field{Key: key, FieldType: FieldTypeString, Value: val}
}

// Stringer creates a Field with a value that implements the fmt.Stringer interface.
func Stringer(key string, val fmt.Stringer) Field {
	return Field{Key: key, FieldType: FieldTypeStringer, Value: val}
}

type _int interface {
	int | int32 | int64
}

// Int creates a Field with an integer value.
func Int[T _int](key string, val T) Field {
	return Field{Key: key, FieldType: FieldTypeInt, Value: int(val)}
}

// Time creates a Field with a time.Time value.
func Time(key string, val time.Time) Field {
	return Field{Key: key, FieldType: FieldTypeTime, Value: val}
}

// Duration creates a Field with a time.Duration value.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, FieldType: FieldTypeDuration, Value: val}
}

// Error creates a Field with an error value, using "error" as the key.
func Error(err error) Field {
	return Field{Key: "error", FieldType: FieldTypeError, Value: err}
}

// Bool creates a Field with a boolean value.
func Bool(key string, val bool) Field {
	return Field{Key: key, FieldType: FieldTypeBool, Value: val}
}

// Any creates a Field with an arbitrary value.
func Any(key string, val interface{}) Field {
	return Field{Key: key, FieldType: FieldTypeAny, Value: val}
}
