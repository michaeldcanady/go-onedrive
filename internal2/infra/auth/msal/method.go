package msal

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
