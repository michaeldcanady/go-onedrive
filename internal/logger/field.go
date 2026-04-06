package logger

// Field represents a key-value pair to be included in a log entry.
// It consists of a descriptive key, a specific [FieldType], and the underlying value.
type Field struct {
	// Key is the name assigned to the field.
	Key string
	// FieldType identifies the data type of the Value.
	FieldType FieldType
	// Value is the actual data associated with the Key.
	Value any
}
