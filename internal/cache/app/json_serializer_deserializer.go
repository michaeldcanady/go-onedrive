package app

import (
	"encoding/json"

	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
)

var _ domaincache.SerializerDeserializer[any] = (*JSONSerializerDeserializer[any])(nil)

type JSONSerializerDeserializer[T any] struct{}

// Deserialize converts the provided json data to the specified type.
func (j *JSONSerializerDeserializer[T]) Deserialize(data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}

// Serialize converts the provided value to json data.
func (j *JSONSerializerDeserializer[T]) Serialize(value T) ([]byte, error) {
	return json.Marshal(value)
}
