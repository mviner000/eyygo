// config/settings_types.go

package config

import "github.com/mviner000/eyymi/eyygo/shared"

// Debuggable interface for debug settings
type Debuggable interface {
	IsDebug() bool
}

// Config interface for accessing configuration
type Config interface {
	GetDatabaseConfig() shared.DatabaseConfig // Use shared.DatabaseConfig
	SetDatabaseConfig(shared.DatabaseConfig)  // Use shared.DatabaseConfig
	IsDebug() bool
	SetDebug(bool)
}

// DatabaseConfig defines the structure for database settings
type DatabaseConfig struct {
	Engine   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	Options  map[string]string
}

// WebSocketConfig defines the structure for WebSocket settings
type WebSocketConfig struct {
	Port string
}

// SettingsStruct defines the structure for application settings
type SettingsStruct struct {
	TemplateBasePath        string
	CurrentWorkingDirectory string
	WebSocket               WebSocketConfig
	AllowedOrigins          string
	CertFile                string
	KeyFile                 string
	LogFile                 string
	Debug                   bool
	TimeZone                string
	InstalledApps           []string
	Database                shared.DatabaseConfig // Use shared.DatabaseConfig
	Environment             string
	IsDevelopment           bool
}
