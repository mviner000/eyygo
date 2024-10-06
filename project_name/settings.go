package project_name

import (
	"os"
	"strings"

	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/shared"
)

var AppSettings SettingsStruct

type WebSocketConfig struct {
	Port string
}

type SettingsStruct struct {
	Database         config.DatabaseConfig
	Debug            bool
	TimeZone         string
	WebSocket        WebSocketConfig
	CertFile         string
	KeyFile          string
	AllowedOrigins   []string
	TemplateBasePath string
	SecretKey        string
}

func LoadSettings() {
	AppSettings = SettingsStruct{
		SecretKey: os.Getenv("SECRET_KEY"),
		Database: config.DatabaseConfig{
			Engine:   os.Getenv("DB_ENGINE"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Debug:    os.Getenv("DEBUG") == "true",
		TimeZone: os.Getenv("TIME_ZONE"),
		WebSocket: WebSocketConfig{
			Port: os.Getenv("WS_PORT"),
		},
		CertFile:         os.Getenv("CERT_FILE"),
		KeyFile:          os.Getenv("KEY_FILE"),
		AllowedOrigins:   strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		TemplateBasePath: os.Getenv("TEMPLATE_BASE_PATH"),
	}

	// Set the secret key in the shared config
	shared.SetSecretKey(os.Getenv("SECRET_KEY"))

	// Log loaded settings
	config.LogStruct("Loaded settings", AppSettings)
}

func (s *SettingsStruct) GetDatabaseConfig() config.DatabaseConfig {
	return s.Database
}

func (s *SettingsStruct) SetDatabaseConfig(dbConfig config.DatabaseConfig) {
	s.Database = dbConfig
}

func (s *SettingsStruct) IsDebug() bool {
	return s.Debug
}

func (s *SettingsStruct) SetDebug(debug bool) {
	s.Debug = debug
}
