package app

import (
	"fmt"

	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

var _ pkgcache.SerializerDeserializer[serialization.Parsable] = (*KiotaJSONSerializerDeserializer[serialization.Parsable])(nil)

type KiotaJSONSerializerDeserializer[T serialization.Parsable] struct {
	factory serialization.ParsableFactory
}

func NewKiotaJSONSerializerDeserializer[T serialization.Parsable](factory serialization.ParsableFactory) *KiotaJSONSerializerDeserializer[T] {
	return &KiotaJSONSerializerDeserializer[T]{
		factory: factory,
	}
}

// Deserialize implements [pkgcache.SerializerDeserializer].
func (j *KiotaJSONSerializerDeserializer[T]) Deserialize(data []byte) (T, error) {
	var v T

	node, err := jsonserialization.NewJsonParseNode(data)
	if err != nil {
		return v, err
	}

	value, err := node.GetObjectValue(j.factory)
	if err != nil {
		return v, err
	}

	v, ok := value.(T)
	if !ok {
		return v, fmt.Errorf("value is not %T", v)
	}

	return v, nil
}

// Serialize implements [pkgcache.SerializerDeserializer].
func (j *KiotaJSONSerializerDeserializer[T]) Serialize(value T) ([]byte, error) {
	writer := jsonserialization.NewJsonSerializationWriter()

	if err := writer.WriteObjectValue("", value); err != nil {
		return nil, err
	}

	return writer.GetSerializedContent()
}
