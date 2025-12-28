package abstractions

// Codec defines how to encode/decode a profile.
type Codec interface {
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
}
