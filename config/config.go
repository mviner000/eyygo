package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/mviner000/eyymi/eyygo/utils"
)

var AppSettings SettingsStruct
var ProjectRoot string

func init() {
	var err error
	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Get the absolute path of the current executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error finding executable path: %v", err)
	}
	log.Printf("Executable path: %s", execPath)

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	log.Printf("Executable directory: %s", execDir)

	// Set ProjectRoot to the parent directory of the executable
	ProjectRoot = filepath.Dir(execDir)
	log.Printf("Project root set to: %s", ProjectRoot)

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", cwd)
	}
}

// LoadSettings loads the application settings, using the provided defaults
func LoadSettings(defaultSettings SettingsStruct) {
	AppSettings = SettingsStruct{
		Environment:    getEnv("NODE_ENV", defaultSettings.Environment),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", defaultSettings.AllowedOrigins),
		CertFile:       getEnv("CERT_FILE", defaultSettings.CertFile),
		KeyFile:        getEnv("KEY_FILE", defaultSettings.KeyFile),
		LogFile:        getEnv("LOG_FILE", defaultSettings.LogFile),
		Debug:          getEnv("DEBUG", fmt.Sprintf("%v", defaultSettings.Debug)) == "true",
		TimeZone:       getEnv("TIME_ZONE", defaultSettings.TimeZone),
		InstalledApps:  defaultSettings.InstalledApps,
		Database: DatabaseConfig{
			Engine:   getEnv("DB_ENGINE", defaultSettings.Database.Engine),
			Name:     getEnv("DB_NAME", defaultSettings.Database.Name),
			User:     getEnv("DB_USER", defaultSettings.Database.User),
			Password: getEnv("DB_PASSWORD", defaultSettings.Database.Password),
			Host:     getEnv("DB_HOST", defaultSettings.Database.Host),
			Port:     getEnv("DB_PORT", defaultSettings.Database.Port),
		},
		WebSocket: WebSocketConfig{
			Port: getEnv("WS_PORT", defaultSettings.WebSocket.Port),
		},
	}

	AppSettings.IsDevelopment = AppSettings.Environment == "development"

	// Print debug status immediately after loading
	fmt.Printf("Debug setting loaded: %v\n", AppSettings.Debug)

	// Log the allowed origins
	log.Printf("ALLOWED_ORIGINS: %s\n", AppSettings.AllowedOrigins)

	// Use DebugLog instead of fmt.Printf for consistency
	DebugLog("Successfully loaded settings from environment variables")
	DebugLog("Database Engine: %s", AppSettings.Database.Engine)
	DebugLog("Database Name: %s", AppSettings.Database.Name)
}

// getEnv retrieves the environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetWebSocketPort() string {
	return AppSettings.WebSocket.Port
}

func GetAllowedOrigins() string {
	return AppSettings.AllowedOrigins
}

func IsDevelopment() bool {
	return AppSettings.Environment == "development" || AppSettings.Debug
}

func GetCertFile() string {
	return AppSettings.CertFile
}

func GetKeyFile() string {
	return AppSettings.KeyFile
}

func GetInstalledApps() []string {
	return AppSettings.InstalledApps
}

func AddInstalledApp(appName string) {
	if !utils.Contains(AppSettings.InstalledApps, appName) {
		AppSettings.InstalledApps = append(AppSettings.InstalledApps, appName)
	}
}
