package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigDir returns the absolute path to ~/.config/button/apps/.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "button", "apps"), nil
}

// EnsureConfigDir creates the config directory if it doesn't already exist.
func EnsureConfigDir() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// ReadApps scans the config directory for *.yaml/*.yml files, parses each one,
// and resolves platform-specific shortcut keys based on the current OS.
// Files that are empty or missing the required "app" field are skipped and
// reported as warnings instead of causing a hard error.
func ReadApps() (AppsResponse, error) {
	dir, err := ConfigDir()
	if err != nil {
		return AppsResponse{}, err
	}

	seen := make(map[string]bool)
	var files []string
	for _, ext := range []string{"*.yaml", "*.yml"} {
		matches, err := filepath.Glob(filepath.Join(dir, ext))
		if err != nil {
			return AppsResponse{}, fmt.Errorf("failed to glob config files: %w", err)
		}
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				files = append(files, m)
			}
		}
	}

	// Sort all files together alphabetically by base name without extension,
	// so .yaml and .yml files for the same app sort alongside each other correctly.
	sort.Slice(files, func(i, j int) bool {
		ni := strings.TrimSuffix(strings.ToLower(filepath.Base(files[i])), filepath.Ext(files[i]))
		nj := strings.TrimSuffix(strings.ToLower(filepath.Base(files[j])), filepath.Ext(files[j]))
		return ni < nj
	})

	var resp AppsResponse
	resp.Apps = make([]AppConfig, 0, len(files))

	for _, file := range files {
		name := filepath.Base(file)

		info, err := os.Stat(file)
		if err != nil {
			resp.Warnings = append(resp.Warnings, fmt.Sprintf("%s: failed to stat — %s", name, err))
			continue
		}

		app, err := parseAppFile(file)
		if err != nil {
			resp.Warnings = append(resp.Warnings, fmt.Sprintf("%s: failed to parse — %s", name, err))
			continue
		}

		if app.App == "" {
			resp.Warnings = append(resp.Warnings, fmt.Sprintf("%s: skipped — file is empty or missing the \"app\" field", name))
			continue
		}

		app.ModTime = info.ModTime().Unix()
		resp.Apps = append(resp.Apps, app)
	}

	return resp, nil
}

// parseAppFile reads and unmarshals a single YAML file into an AppConfig.
// Returns a zero-value AppConfig (empty App field) if the file is empty or whitespace-only.
func parseAppFile(path string) (AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return AppConfig{}, nil
	}

	var app AppConfig
	if err := yaml.Unmarshal(data, &app); err != nil {
		return AppConfig{}, err
	}

	return app, nil
}

// resolveKeys walks every shortcut in the app and collapses the
// platform-specific fields into the single Keys field based on runtime.GOOS.
// Priority: linux/macos override > keys fallback.
func resolveKeys(app *AppConfig) {
	for i := range app.Groups {
		for j := range app.Groups[i].Shortcuts {
			s := &app.Groups[i].Shortcuts[j]

			switch runtime.GOOS {
			case "linux":
				if len(s.Linux) > 0 {
					s.Keys = s.Linux
				}
			case "darwin":
				if len(s.MacOS) > 0 {
					s.Keys = s.MacOS
				}
			}
			// If neither override was set, Keys keeps its original value (the fallback).
		}
	}
}
