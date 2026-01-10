package config

import "github.com/spf13/viper"

type ViperAdapter struct {
	base *viper.Viper
}

func NewViperAdapter(v *viper.Viper) *ViperAdapter {
	return &ViperAdapter{
		base: v,
	}
}

func (v *ViperAdapter) SetConfigFile(path string) {
	v.base.SetConfigFile(path)
}

func (v *ViperAdapter) ReadInConfig() error {
	return v.base.ReadInConfig()
}

func (v *ViperAdapter) GetString(key string) string {
	return v.base.GetString(key)
}

func (v *ViperAdapter) Get(key string) interface{} {
	return v.base.Get(key)
}

func (v *ViperAdapter) Set(key string, value interface{}) {
	v.base.Set(key, value)
}

func (v *ViperAdapter) WriteConfig() error {
	return v.base.WriteConfig()
}
