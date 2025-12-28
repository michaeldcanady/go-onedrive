package jsoncodec

import "encoding/json"

type JSONCodec struct{}

func New() *JSONCodec {
	return &JSONCodec{}
}

func (c *JSONCodec) Encode(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func (c *JSONCodec) Decode(data []byte, v any) error {
	// treat empty files as "no data"
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}
