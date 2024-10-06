package shared

import "sync"

var (
	config     *Config
	configOnce sync.Once
)

type Config struct {
	SecretKey string
}

func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{}
	})
	return config
}

func SetSecretKey(key string) {
	GetConfig().SecretKey = key
}
