package conf

import (
	"github.com/mviner000/eyymi/project_name"
)

var Settings *project_name.SettingsStruct

func init() {
	project_name.LoadSettings()
	Settings = &project_name.AppSettings
}

// GetSettings returns the global settings
func GetSettings() *project_name.SettingsStruct {
	return Settings
}
