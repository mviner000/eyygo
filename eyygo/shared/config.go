// File: shared/config.go

package shared

import "sync"

var (
	config     *Config
	configOnce sync.Once
)

type DatabaseConfig struct {
	Engine   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

type Config struct {
	SecretKey string
	Database  DatabaseConfig
	Debug     bool
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

func SetDatabaseConfig(dbConfig DatabaseConfig) {
	GetConfig().Database = dbConfig
}

func SetDebug(debug bool) {
	GetConfig().Debug = debug
}
