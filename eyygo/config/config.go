package config

func GetWebSocketPort(settings SettingsStruct) string {
	return settings.WebSocket.Port
}

func GetAllowedOrigins(settings SettingsStruct) string {
	return settings.AllowedOrigins
}

func GetCertFile(settings SettingsStruct) string {
	return settings.CertFile
}

func GetKeyFile(settings SettingsStruct) string {
	return settings.KeyFile
}

func GetInstalledApps(settings SettingsStruct) []string {
	return settings.InstalledApps
}

func AddInstalledApp(settings *SettingsStruct, appName string) {
	if !contains(settings.InstalledApps, appName) {
		settings.InstalledApps = append(settings.InstalledApps, appName)
	}
}

func IsDevelopment(settings Debuggable) bool {
	return settings.IsDebug()
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
