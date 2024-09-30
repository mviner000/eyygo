package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Settings struct {
	Environment    string   `json:"environment"`
	WebSocketPort  string   `json:"webSocketPort"`
	AllowedOrigins string   `json:"allowedOrigins"`
	CertFile       string   `json:"certFile"`
	KeyFile        string   `json:"keyFile"`
	LogFile        string   `json:"logFile"`
	IsDevelopment  bool     `json:"isDevelopment"`
	InstalledApps  []string `json:"installedApps"`
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
	}

	AppSettings.IsDevelopment = AppSettings.Environment == "development"
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
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "username")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "database")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
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