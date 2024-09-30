package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Engine   string `json:"engine"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
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

const settingsFile = "config.json"

func init() {
	loadSettings()
}

func loadSettings() {
	file, err := os.Open(settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Debug: config.json not found. Using default settings.")
			AppSettings = getDefaultSettings()
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
		log.Printf("Debug: Successfully loaded settings from config.json")
	}

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
	case "postgresql":
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			db.User, db.Password, db.Host, db.Port, db.Name)
	case "mysql":
		dbURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			db.User, db.Password, db.Host, db.Port, db.Name)
	case "sqlite3":
		dbPath := filepath.Join(getEnv("BASE_DIR", "."), db.Name)
		dbURL = fmt.Sprintf("sqlite3://%s", filepath.ToSlash(dbPath))
	default:
		log.Printf("Debug: Unsupported database engine: %s, falling back to SQLite", db.Engine)
		dbPath := filepath.Join(getEnv("BASE_DIR", "."), "db.sqlite3")
		dbURL = fmt.Sprintf("sqlite3://%s", filepath.ToSlash(dbPath))
	}
	log.Printf("Debug: Database URL: %s", dbURL)
	return dbURL
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