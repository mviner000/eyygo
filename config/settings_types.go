package config

type Debuggable interface {
	IsDebug() bool
}

type Config interface {
	GetDatabaseConfig() DatabaseConfig
	SetDatabaseConfig(DatabaseConfig)
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
	Database                DatabaseConfig
	Environment             string
	IsDevelopment           bool
}
