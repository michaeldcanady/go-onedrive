package cacheservice

import (
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
)

var _ abstractions.SerializerDeserializer[any] = (*JSONSerializerDeserializer[any])(nil)

type JSONSerializerDeserializer[T any] struct{}

// Deserialize implements [abstractions.SerializerDeserializer].
func (j *JSONSerializerDeserializer[T]) Deserialize(data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}

// Serialize implements [abstractions.SerializerDeserializer].
func (j *JSONSerializerDeserializer[T]) Serialize(value T) ([]byte, error) {
	return json.Marshal(value)
}
