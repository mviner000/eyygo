package config

// DatabaseConfig defines the structure for database settings
type DatabaseConfig struct {
	// Engine is the database engine (e.g., sqlite3, postgres, mysql)
	Engine string

	// Name is the name of the database
	Name string

	// User is the username for the database
	User string

	// Password is the password for the database
	Password string

	// Host is the host for the database
	Host string

	// Port is the port for the database
	Port string
}

// WebSocketConfig defines the structure for WebSocket settings
type WebSocketConfig struct {
	// Port is the port for the WebSocket
	Port string
}

// SettingsStruct defines the structure for application settings
type SettingsStruct struct {
	// TemplateBasePath is the base path for templates
	TemplateBasePath string

	// CurrentWorkingDirectory is the current working directory of the application
	CurrentWorkingDirectory string

	// Database settings
	Database DatabaseConfig

	// WebSocket settings
	WebSocket WebSocketConfig

	// AllowedOrigins is the list of allowed origins for the application
	AllowedOrigins string

	// CertFile is the path to the certificate file
	CertFile string

	// KeyFile is the path to the key file
	KeyFile string

	// LogFile is the path to the log file
	LogFile string

	// Debug is the debug setting for the application
	Debug bool
}
