package logger

// FieldType is an enumeration of the supported data types for log fields.
type FieldType int

const (
	// FieldTypeString represents a string-valued log field.
	FieldTypeString FieldType = iota
	// FieldTypeInt represents an integer-valued log field.
	FieldTypeInt
	// FieldTypeTime represents a time-valued log field.
	FieldTypeTime
	// FieldTypeDuration represents a duration-valued log field.
	FieldTypeDuration
	// FieldTypeError represents an error-valued log field.
	FieldTypeError
	// FieldTypeBool represents a boolean-valued log field.
	FieldTypeBool
)
