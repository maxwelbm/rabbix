package sett

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SettItf interface {
	LoadSettings() map[string]string
	SaveSettings(settings map[string]string)
}

var _ SettItf = (*Sett)(nil)

type Sett struct {
	path string
}

func New() *Sett {
	return &Sett{
		path: getSettingsPath(),
	}
}

func getSettingsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix", "settings.json")
}

func (s *Sett) LoadSettings() map[string]string {
	settings := map[string]string{}

	if data, err := os.ReadFile(s.path); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	return settings
}

func (s *Sett) SaveSettings(settings map[string]string) {
	_ = os.MkdirAll(filepath.Dir(s.path), os.ModePerm)

	data, _ := json.MarshalIndent(settings, "", "  ")
	_ = os.WriteFile(s.path, data, 0644)
}
