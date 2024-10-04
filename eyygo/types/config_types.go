package types

type DatabaseConfig struct {
	Engine   string            `json:"engine"`
	Name     string            `json:"name"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Host     string            `json:"host"`
	Port     string            `json:"port"`
	Options  map[string]string `json:"options"`
}

type Settings struct {
	Environment    string         `json:"environment"`
	WebSocketPort  string         `json:"webSocketPort"`
	AllowedOrigins string         `json:"allowedOrigins"`
	CertFile       string         `json:"certFile"`
	KeyFile        string         `json:"keyFile"`
	LogFile        string         `json:"logFile"`
	IsDevelopment  bool           `json:"isDevelopment"`
	Debug          bool           `json:"debug"`
	InstalledApps  []string       `json:"installedApps"`
	Database       DatabaseConfig `json:"database"`
}
