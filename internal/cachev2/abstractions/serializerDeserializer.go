package abstractions

type SerializerDeserializer[T any] interface {
	Serializer[T]
	Deserializer[T]
}
