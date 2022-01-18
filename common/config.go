package common

var AppConfig *Config = nil

func init() {
	AppConfig = GetDefaultConfig()
}
