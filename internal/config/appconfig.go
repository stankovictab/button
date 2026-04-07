package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// UserConfig holds persistent user preferences stored in config.yaml.
// New fields can be added here as Button gains more settings.
type UserConfig struct {
	HasSeenWelcome bool   `yaml:"hasSeenWelcome" json:"hasSeenWelcome"`
	LastSortMode   string `yaml:"lastSortMode,omitempty" json:"lastSortMode"`
	GroupByTag     bool   `yaml:"groupByTag" json:"groupByTag"`
}

// BaseDir returns the base button config directory (parent of the apps dir).
// On Windows: %LOCALAPPDATA%\button\, elsewhere: ~/.config/button/.
func BaseDir() (string, error) {
	if runtime.GOOS == "windows" {
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not determine home directory: %w", err)
			}
			local = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(local, "button"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "button"), nil
}

// ReadUserConfig reads config.yaml from the base config directory.
// Returns a zero-value UserConfig if the file doesn't exist.
func ReadUserConfig() (UserConfig, error) {
	dir, err := BaseDir()
	if err != nil {
		return UserConfig{}, err
	}

	data, err := os.ReadFile(filepath.Join(dir, "config.yaml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return UserConfig{}, nil
		}
		return UserConfig{}, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	var cfg UserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return UserConfig{}, fmt.Errorf("failed to parse config.yaml: %w", err)
	}
	return cfg, nil
}

// WriteUserConfig writes the given config to config.yaml in the base config directory.
func WriteUserConfig(cfg UserConfig) error {
	dir, err := BaseDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0644)
}
