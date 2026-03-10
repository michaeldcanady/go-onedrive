package cache

import (
	"encoding/json"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
)

var _ domaincache.SerializerDeserializer[any] = (*JSONSerializerDeserializer[any])(nil)

type JSONSerializerDeserializer[T any] struct{}

// Deserialize implements [domaincache.SerializerDeserializer].
func (j *JSONSerializerDeserializer[T]) Deserialize(data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}

// Serialize implements [domaincache.SerializerDeserializer].
func (j *JSONSerializerDeserializer[T]) Serialize(value T) ([]byte, error) {
	return json.Marshal(value)
}
