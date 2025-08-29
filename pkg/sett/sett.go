package sett

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SettItf interface {
	LoadSettings() map[string]string
	SaveSettings(settings map[string]string)
	GetBaseDir() string
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
	baseDir := getBaseDir()
	pathSett := filepath.Join(baseDir, "settings.json")

	// Guarantees base directory
	_ = os.MkdirAll(baseDir, os.ModePerm)

	// Defaults
	defaultSett := map[string]string{"sett": "local.json"}
	defaultCfg := map[string]string{
		"auth":       "Z3Vlc3Q6Z3Vlc3Q=",
		"host":       "http://localhost:15672",
		"output_dir": baseDir,
	}

	// Create settings.json if it does not exist.
	if _, err := os.Stat(pathSett); os.IsNotExist(err) {
		if data, err := json.MarshalIndent(defaultSett, "", "  "); err == nil {
			_ = os.WriteFile(pathSett, data, 0644)
		}
	}

	// Loads settings.json
	settings := map[string]string{}
	if data, err := os.ReadFile(pathSett); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	// Guarantees the "sett" key and persists if necessary.
	if settings["sett"] == "" {
		settings["sett"] = defaultSett["sett"]
		if data, err := json.MarshalIndent(settings, "", "  "); err == nil {
			_ = os.WriteFile(pathSett, data, 0644)
		}
	}

	// Target configuration file path
	targetPath := filepath.Join(baseDir, settings["sett"])

	// Creates the target file with defaults if it doesn't exist.
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if data, err := json.MarshalIndent(defaultCfg, "", "  "); err == nil {
			_ = os.WriteFile(targetPath, data, 0644)
		}
	}

	return targetPath
}

func getBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix")
}

func (s *Sett) GetBaseDir() string {
	return getBaseDir()
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
