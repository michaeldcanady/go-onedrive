package abstractions

type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
}
