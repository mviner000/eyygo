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
	SecretKey      string
	Database       DatabaseConfig
	Debug          bool
	AllowedOrigins []string
	Environment    string
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

func GetSecretKey() string {
	return GetConfig().SecretKey
}

func SetDatabaseConfig(dbConfig DatabaseConfig) {
	GetConfig().Database = dbConfig
}

func SetDebug(debug bool) {
	GetConfig().Debug = debug
}

func SetAllowedOrigins(origins []string) {
	GetConfig().AllowedOrigins = origins
}

func SetEnvironment(env string) {
	GetConfig().Environment = env
}
