package logging

import (
	"context"
	"time"
)

// Field represents a key-value pair to be logged, along with its type.
type Field struct {
	// Key is the name of the field.
	Key string
	// FieldType indicates the type of the field.
	FieldType FieldType
	// Value holds the actual Value of the field.
	Value any
}

// String creates a Field with a string value.
func String(key, val string) Field {
	return Field{Key: key, FieldType: FieldTypeString, Value: val}
}

type _int interface {
	int | int32 | int64
}

// Int creates a Field with an integer value.
func Int[T _int](key string, val T) Field {
	return Field{Key: key, FieldType: FieldTypeInt, Value: int(val)}
}

// Any creates a Field with an arbitrary value.
func Any(key string, val any) Field {
	return Field{Key: key, FieldType: FieldTypeAny, Value: val}
}

// Bool creates a Field with a boolean value.
func Bool(key string, val bool) Field {
	return Field{Key: key, FieldType: FieldTypeBool, Value: val}
}

// Duration creates a Field with a boolean value.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, FieldType: FieldTypeDuration, Value: val}
}

func Strings(key string, val []string) Field {
	return Field{Key: key, FieldType: FieldTypeStrings, Value: val}
}

func Time(key string, val time.Time) Field {
	return Field{Key: key, FieldType: FieldTypeTime, Value: val}
}

func Error(val error) Field {
	return Field{Key: "error", FieldType: FieldTypeError, Value: val}
}

type ctxKey struct{}

type Fields map[string]any

func WithFields(ctx context.Context, fields ...Field) context.Context {
	existing := FromContextFields(ctx)
	merged := make(Fields, len(existing)+len(fields))

	for k, v := range existing {
		merged[k] = v
	}
	for _, f := range fields {
		merged[f.Key] = f.Value
	}

	return context.WithValue(ctx, ctxKey{}, merged)
}

func FromContextFields(ctx context.Context) Fields {
	if v := ctx.Value(ctxKey{}); v != nil {
		if f, ok := v.(Fields); ok {
			return f
		}
	}
	return Fields{}
}
