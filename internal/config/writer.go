package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// appNameToFilename converts an app display name to a YAML filename.
// e.g. "My Cool App" → "my-cool-app.yaml"
func appNameToFilename(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	dashed := strings.ReplaceAll(lower, " ", "-")
	return dashed + ".yaml"
}

// yamlAppConfig is the subset of AppConfig that gets persisted to YAML
// (excludes runtime-only fields like ModTime).
type yamlAppConfig struct {
	App    string  `yaml:"app"`
	Icon   string  `yaml:"icon"`
	Groups []Group `yaml:"groups"`
}

// marshalApp produces clean YAML bytes for the given AppConfig.
func marshalApp(app AppConfig) ([]byte, error) {
	out := yamlAppConfig{
		App:    app.App,
		Icon:   app.Icon,
		Groups: app.Groups,
	}
	return yaml.Marshal(out)
}

// CreateApp writes a new YAML file for the given app config.
// Returns an error if a file for this app name already exists.
func CreateApp(app AppConfig) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	filename := appNameToFilename(app.App)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("app file %q already exists", filename)
	}

	data, err := marshalApp(app)
	if err != nil {
		return fmt.Errorf("failed to marshal app config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// UpdateApp replaces the YAML file for an app. If the app name changed,
// the old file is removed and a new one is created.
// Returns a non-empty warning string if the new filename would overwrite
// a different existing file — in that case the caller should confirm first
// by calling UpdateApp again with force=true (via UpdateAppForce).
func UpdateApp(oldAppName string, app AppConfig, force bool) (warning string, err error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	oldFilename := appNameToFilename(oldAppName)
	newFilename := appNameToFilename(app.App)
	oldPath := filepath.Join(dir, oldFilename)
	newPath := filepath.Join(dir, newFilename)

	// If name changed and new file already exists (and it's not the same file), warn.
	if oldFilename != newFilename {
		if _, err := os.Stat(newPath); err == nil && !force {
			return fmt.Sprintf(
				"A file named %q already exists. Saving will overwrite it.",
				newFilename,
			), nil
		}
	}

	data, err := marshalApp(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal app config: %w", err)
	}

	// Write new file first, then remove old if different.
	if err := os.WriteFile(newPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write app config: %w", err)
	}

	if oldFilename != newFilename {
		// Remove old file (ignore error if it doesn't exist).
		os.Remove(oldPath)
	}

	return "", nil
}

// DeleteApp removes the YAML file for the given app name.
func DeleteApp(appName string) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	filename := appNameToFilename(appName)
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("app file %q does not exist", filename)
	}

	return os.Remove(path)
}
