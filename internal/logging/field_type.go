package logging

// FieldType indicates the type of a logging field.
type FieldType int

const (
	// FieldTypeString represents a field with a string value.
	FieldTypeString FieldType = iota
	// FieldTypeInt represents a field with an integer value.
	FieldTypeInt
	// FieldTypeAny represents a field with an arbitrary value.
	FieldTypeAny
)
