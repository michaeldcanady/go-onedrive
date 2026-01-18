package abstractions

type Deserializer[T any] interface {
	Deserialize([]byte) (T, error)
}
