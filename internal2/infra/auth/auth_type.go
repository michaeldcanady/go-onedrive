package auth

const (
	InteractiveBrowserAuthType = "interactiveBrowser"
)

const (
	AuthTypeUnknownString            = "unknown"
	AuthTypeInteractiveBrowserString = "interactiveBrowser"
	AuthTypeDeviceCodeString         = "deviceCode"
	AuthTypeROPCString               = "ropc"
	AuthTypeClientSecretString       = "clientSecret"
)

type AuthType int32

func ParseAuthType(str string) AuthType {
	authType, ok := map[string]AuthType{
		AuthTypeUnknownString:            AuthTypeUnknown,
		AuthTypeInteractiveBrowserString: AuthTypeInteractiveBrowser,
		AuthTypeDeviceCodeString:         AuthTypeDeviceCode,
		AuthTypeROPCString:               AuthTypeROPC,
		AuthTypeClientSecretString:       AuthTypeClientSecret,
	}[str]
	if !ok {
		return AuthTypeUnknown
	}
	return authType
}

func (a *AuthType) UnmarshalJSON(data []byte) error {
	str := string(data)
	*a = ParseAuthType(str)
	return nil
}

func (a AuthType) MarshalJSON() ([]byte, error) {
	str := a.String()
	return []byte(str), nil
}

func (a *AuthType) UnmarshalYAML(data []byte) error {
	str := string(data)
	*a = ParseAuthType(str)
	return nil
}

func (a AuthType) MarshalYAML() (interface{}, error) {
	str := a.String()
	return str, nil
}

func (a AuthType) String() string {
	str, ok := map[AuthType]string{
		AuthTypeUnknown:            AuthTypeUnknownString,
		AuthTypeInteractiveBrowser: AuthTypeInteractiveBrowserString,
		AuthTypeDeviceCode:         AuthTypeDeviceCodeString,
		AuthTypeROPC:               AuthTypeROPCString,
		AuthTypeClientSecret:       AuthTypeClientSecretString,
	}[a]
	if !ok {
		return AuthTypeUnknown.String()
	}
	return str
}

const (
	AuthTypeUnknown AuthType = iota - 1
	AuthTypeInteractiveBrowser
	AuthTypeDeviceCode
	AuthTypeROPC
	AuthTypeClientSecret
)
