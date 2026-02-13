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

func (a AuthType) String() string {

}

const (
	AuthTypeUnknown AuthType = iota - 1
	AuthTypeInteractiveBrowser
	AuthTypeDeviceCode
	AuthTypeROPC
	AuthTypeClientSecret
)
