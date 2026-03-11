package domain

type FieldType int

const (
	FieldTypeUnknown FieldType = iota
	FieldTypeString
	FieldTypeInt
	FieldTypeAny
	FieldTypeBool
	FieldTypeDuration
	FieldTypeStrings
	FieldTypeTime
	FieldTypeError
)
