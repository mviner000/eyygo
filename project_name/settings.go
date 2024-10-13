package project_name

import (
	"os"
	"path/filepath"

	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/shared"
	"github.com/mviner000/eyymi/eyygo/utils"
)

var AppSettings SettingsStruct

type WebSocketConfig struct {
	Port string
}

type SettingsStruct struct {
	Database         shared.DatabaseConfig
	Debug            bool
	TimeZone         string
	WebSocket        WebSocketConfig
	CertFile         string
	KeyFile          string
	AllowedOrigins   []string
	TemplateBasePath string
	SecretKey        string
	LogFile          string
	InstalledApps    []string
	Environment      string
	IsDevelopment    bool
}

// Helper function to create app paths
func createAppPaths(apps []string) []string {
	return append([]string{}, apps...)
}

// LoadSettings initializes application settings
func LoadSettings() {
	dbConfig := shared.DatabaseConfig{
		Engine:   os.Getenv("DB_ENGINE"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}
	debug := os.Getenv("DEBUG") == "true"

	// Initialize the shared config
	shared.SetSecretKey(os.Getenv("SECRET_KEY"))
	shared.SetDatabaseConfig(dbConfig)
	shared.SetDebug(debug)

	// Get project root
	projectRoot := utils.GetProjectRoot(debug)

	// Define INSTALLED_APPS with the new structure
	installedApps := createAppPaths([]string{
		"eyygo.admin",
		"eyygo.sessions",
		"eyygo.auth",
		"eyygo.contenttypes",
		"project_name.posts",
		"project_name.notes",
	})

	// Initialize the settings struct
	AppSettings = SettingsStruct{
		TemplateBasePath: filepath.Join(projectRoot, os.Getenv("TEMPLATE_BASE_PATH")),
		InstalledApps:    installedApps,
		WebSocket:        WebSocketConfig{Port: os.Getenv("WS_PORT")},
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		CertFile:         filepath.Join(projectRoot, os.Getenv("CERT_FILE")),
		KeyFile:          filepath.Join(projectRoot, os.Getenv("KEY_FILE")),
		LogFile:          filepath.Join(projectRoot, os.Getenv("LOG_FILE")),
		Debug:            debug,
		TimeZone:         os.Getenv("TIME_ZONE"),
		Database:         dbConfig,
		Environment:      os.Getenv("ENVIRONMENT"),
		IsDevelopment:    os.Getenv("ENVIRONMENT") == "development",
	}

	// Log loaded settings
	config.LogStruct("Loaded settings", AppSettings)
}

func (s *SettingsStruct) GetDatabaseConfig() shared.DatabaseConfig {
	return s.Database
}

func (s *SettingsStruct) SetDatabaseConfig(dbConfig shared.DatabaseConfig) {
	s.Database = dbConfig
}

func (s *SettingsStruct) IsDebug() bool {
	return s.Debug
}

func (s *SettingsStruct) SetDebug(debug bool) {
	s.Debug = debug
}
