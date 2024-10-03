package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mviner000/eyymi/types"
	"github.com/mviner000/eyymi/utils"
)

var AppSettings types.Settings
var ProjectRoot string

const settingsFile = "config/config.json"

func init() {
	var err error
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

	loadSettings() // Load settings first
	initLogger()   // Then initialize logger

	if AppSettings.Debug {
		DebugLog("Debug: %v", AppSettings.Debug)
		DebugLog("Project root: %s", ProjectRoot)
	}
}

func loadSettings() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	configPath := filepath.Join(cwd, settingsFile)

	file, err := os.Open(configPath)
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

	// Print debug status immediately after loading
	fmt.Printf("Debug setting loaded: %v\n", AppSettings.Debug)

	// Use DebugLog instead of fmt.Printf for consistency
	DebugLog("Successfully loaded settings from %s", configPath)
	DebugLog("Database Engine: %s", AppSettings.Database.Engine)
	DebugLog("Database Name: %s", AppSettings.Database.Name)
}

func getDefaultSettings() types.Settings {
	getEnv := utils.GetEnv

	return types.Settings{
		Environment:    getEnv("NODE_ENV", "development"),
		WebSocketPort:  getEnv("WS_PORT", "3000"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "https://eyymi.site"),
		CertFile:       getEnv("CERT_FILE", ""),
		KeyFile:        getEnv("KEY_FILE", ""),
		LogFile:        getEnv("LOG_FILE", "server.log"),
		Debug:          getEnv("DEBUG", "false") == "true",
		InstalledApps:  []string{},
		Database: types.DatabaseConfig{
			Engine:   getEnv("DB_ENGINE", "sqlite3"),
			Name:     getEnv("DB_NAME", "db.sqlite3"),
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

func GetInstalledApps() []string {
	return AppSettings.InstalledApps
}

func AddInstalledApp(appName string) {
	if !utils.Contains(AppSettings.InstalledApps, appName) {
		AppSettings.InstalledApps = append(AppSettings.InstalledApps, appName)
		saveSettings()
	}
}
