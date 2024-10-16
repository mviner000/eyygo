package project_name

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mviner000/eyymi/src/config"
	"github.com/mviner000/eyymi/src/shared"
	"github.com/mviner000/eyymi/src/utils"
)

const (
	FullProjectName = "project_name"
)

var AppSettings SettingsStruct

type WebSocketConfig struct {
	Port string
}

type CSRFConfig struct {
	Secret    string
	TokenName string
	Secure    bool
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
	CSRF             CSRFConfig
	FullProjectName  string
}

// Helper function to create app paths
func createAppPaths(apps []string) []string {
	return append([]string{}, apps...)
}

// Logger function to print bold green text
func logSuccess(message string) {
	fmt.Printf("\033[1;32m%s\033[0m\n", message) // \033[1;32m is for bold green
}

// LoadSettings initializes application settings
func LoadSettings() {
	// Load environment variables from the .env file using godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set Environment
	environment := os.Getenv("NODE_ENV")
	shared.SetEnvironment(environment)

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
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set")
	}

	// Log the loaded secret key
	log.Printf("Loaded Secret Key: %s\n", os.Getenv("SECRET_KEY"))

	// Force log the SECRET_KEY for debugging
	fmt.Printf("\033[1;31m[DEBUG] SECRET_KEY: %s\033[0m\n", secretKey) // Prints in bold red

	shared.SetSecretKey(secretKey)
	logSuccess("SECRET_KEY successfully loaded")

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

	// Set AllowedOrigins
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	shared.SetAllowedOrigins(allowedOrigins)

	// Initialize the settings struct
	AppSettings = SettingsStruct{
		FullProjectName:  FullProjectName,
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
		CSRF: CSRFConfig{
			Secret:    os.Getenv("CSRF_SECRET"),
			TokenName: os.Getenv("CSRF_TOKEN_NAME"),
			Secure:    os.Getenv("CSRF_SECURE") == "true",
		},
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
