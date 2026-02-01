package config

type Configuration2 interface {
	Get(string) interface{}
	GetString(string) string
	Set(string, interface{})
	SetConfigFile(path string)
	ReadInConfig() error
	WriteConfig() error
}
