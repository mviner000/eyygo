package app

import (
	"github.com/mviner000/eyymi/config"
)

// Settings holds the configuration for the application
var Settings = config.SettingsStruct{
	TemplateBasePath:        "eyygo/",
	CurrentWorkingDirectory: "",
	Database: config.DatabaseConfig{
		Engine:   "sqlite3",
		Name:     "db.sqlite3",
		User:     "",
		Password: "",
		Host:     "",
		Port:     "",
	},
	WebSocket: config.WebSocketConfig{
		Port: "3000",
	},
	AllowedOrigins: "https://eyymi.site",
	CertFile:       "",
	KeyFile:        "",
	LogFile:        "server.log",
	Debug:          false,
}
