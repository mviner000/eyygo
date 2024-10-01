package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

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
	InstalledApps  []string       `json:"installedApps"`
	Database       DatabaseConfig `json:"database"`
}

var AppSettings Settings
var ProjectRoot string

const settingsFile = "config/config.json"

func init() {
	var err error
	ProjectRoot, err = filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	log.Printf("Debug: Project root: %s", ProjectRoot)

	// Load settings
	loadSettings()
}

func loadSettings() {
	// Get the current working directory (in case ProjectRoot is not absolute)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Construct the absolute path to the config file
	configPath := filepath.Join(cwd, settingsFile)

	log.Printf("Debug: Attempting to load config from: %s", configPath)

	// Try to open the config file
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Debug: config.json not found at %s. Using default settings.", configPath)
			AppSettings = getDefaultSettings()
			// Optionally save the default settings if desired
			saveSettings()
		} else {
			log.Fatalf("Error opening config file: %v", err)
		}
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&AppSettings)
		if err != nil {
			log.Fatalf("Error decoding config file: %v", err)
		}
		log.Printf("Debug: Successfully loaded settings from %s", configPath)
	}

	// Set additional properties based on the config
	AppSettings.IsDevelopment = AppSettings.Environment == "development"
	log.Printf("Debug: Database Engine: %s", AppSettings.Database.Engine)
	log.Printf("Debug: Database Name: %s", AppSettings.Database.Name)
}

func getDefaultSettings() Settings {
	return Settings{
		Environment:    getEnv("NODE_ENV", "development"),
		WebSocketPort:  getEnv("WS_PORT", "3000"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "https://eyymi.site"),
		CertFile:       getEnv("CERT_FILE", ""),
		KeyFile:        getEnv("KEY_FILE", ""),
		LogFile:        getEnv("LOG_FILE", "server.log"),
		InstalledApps:  []string{},
		Database: DatabaseConfig{
			Engine:   getEnv("DB_ENGINE", "sqlite3"),  // Changed default to sqlite3
			Name:     getEnv("DB_NAME", "db.sqlite3"), // Changed default name
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			Options:  make(map[string]string),
		},
	}
}

func saveSettings() {
	file, err := os.Create(settingsFile)
	if err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(AppSettings)
	if err != nil {
		log.Fatalf("Error encoding config file: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetLogger() *log.Logger {
	logFile, err := os.OpenFile(AppSettings.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func GetWebSocketPort() string {
	return AppSettings.WebSocketPort
}

func GetAllowedOrigins() string {
	return AppSettings.AllowedOrigins
}

func IsDevelopment() bool {
	return AppSettings.IsDevelopment
}

func GetCertFile() string {
	return AppSettings.CertFile
}

func GetKeyFile() string {
	return AppSettings.KeyFile
}

func GetDatabaseURL() string {
	db := AppSettings.Database
	var dbURL string
	switch db.Engine {
	case "sqlite3":
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current working directory: %v", err)
			cwd = "."
		}
		dbPath := filepath.Join(cwd, db.Name)
		dbURL = dbPath // Ent expects the file path for SQLite, not a URL
	// ... (keep other database cases)
	default:
		log.Printf("Debug: Unsupported database engine: %s, falling back to SQLite", db.Engine)
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting current working directory: %v", err)
			cwd = "."
		}
		dbPath := filepath.Join(cwd, "db.sqlite3")
		dbURL = dbPath
	}
	log.Printf("Debug: Database URL: %s", dbURL)
	return dbURL
}

func EnsureDatabaseExists() error {
	db := AppSettings.Database
	if db.Engine == "sqlite3" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %v", err)
		}
		dbPath := filepath.Join(cwd, db.Name)

		// Check if the file exists
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			// Create the file if it doesn't exist
			file, err := os.Create(dbPath)
			if err != nil {
				return fmt.Errorf("error creating SQLite database file: %v", err)
			}
			file.Close()
			log.Printf("Created SQLite database file: %s", dbPath)
		}
	}
	return nil
}

func GetInstalledApps() []string {
	return AppSettings.InstalledApps
}

func AddInstalledApp(appName string) {
	if !contains(AppSettings.InstalledApps, appName) {
		AppSettings.InstalledApps = append(AppSettings.InstalledApps, appName)
		saveSettings()
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
