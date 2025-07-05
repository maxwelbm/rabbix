package sett

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix", "settings.json")
}

func LoadSettings() map[string]string {
	path := GetSettingsPath()
	settings := map[string]string{}

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	return settings
}

func SaveSettings(settings map[string]string) {
	path := GetSettingsPath()
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

	data, _ := json.MarshalIndent(settings, "", "  ")
	_ = os.WriteFile(path, data, 0644)
}
