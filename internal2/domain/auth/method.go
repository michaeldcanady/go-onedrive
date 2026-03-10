package auth

import "encoding/json"

type Method int32

func (m Method) String() string {
	str, ok := map[Method]string{
		MethodUnknown:            MethodUnknownString,
		MethodInteractiveBrowser: MethodInteractiveBrowserString,
	}[m]
	if !ok {
		return MethodUnknown.String()
	}
	return str
}

func ParseMethod(str string) Method {
	method, ok := map[string]Method{
		MethodUnknownString:            MethodUnknown,
		MethodInteractiveBrowserString: MethodInteractiveBrowser,
	}[str]
	if !ok {
		return MethodUnknown
	}
	return method
}

const (
	MethodUnknownString            = "unknown"
	MethodInteractiveBrowserString = "interactiveBrowser"

	MethodUnknown Method = iota - 1
	MethodInteractiveBrowser
)

func (m Method) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}

func (m *Method) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	*m = ParseMethod(s)
	return nil
}

func (m Method) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Method) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*m = ParseMethod(s)
	return nil
}
