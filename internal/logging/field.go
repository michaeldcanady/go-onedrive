package logging

// Field represents a key-value pair to be logged, along with its type.
type Field struct {
	// key is the name of the field.
	key string
	// fieldType indicates the type of the field.
	fieldType FieldType
	// value holds the actual value of the field.
	value any
}

// String creates a Field with a string value.
func String(key, val string) Field {
	return Field{key: key, fieldType: FieldTypeString, value: val}
}

// Int creates a Field with an integer value.
func Int(key string, val int) Field {
	return Field{key: key, fieldType: FieldTypeInt, value: val}
}

// Any creates a Field with an arbitrary value.
func Any(key string, val any) Field {
	return Field{key: key, fieldType: FieldTypeAny, value: val}
}

// Bool creates a Field with a boolean value.
func Bool(key string, val bool) Field {
	return Field{key: key, fieldType: FieldTypeBool, value: val}
}
