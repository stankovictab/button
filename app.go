package main

import (
	"button/internal/config"
	"button/internal/registry"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// MigrationResult holds the outcome of the .yml → .yaml migration at startup.
type MigrationResult struct {
	Migrated int      `json:"migrated"`
	Warnings []string `json:"warnings"`
}

// AppInfo holds the user-facing app metadata surfaced to the frontend.
type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// App struct
type App struct {
	ctx             context.Context
	migrationResult MigrationResult
	registry        *registry.Registry
}

// NewApp creates a new App application struct
func NewApp(reg *registry.Registry) *App {
	return &App{registry: reg}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Ensure the config directory exists before any reads
	if err := config.EnsureConfigDir(); err != nil {
		fmt.Println("Warning: could not create config directory:", err)
	}

	// Seed button.yaml if the apps directory is empty
	if err := config.EnsureDefaultApp(); err != nil {
		fmt.Println("Warning: could not create default app file:", err)
	}

	// Migrate .yml files to .yaml
	migrated, migrationWarnings := config.MigrateYMLToYAML()
	a.migrationResult = MigrationResult{
		Migrated: migrated,
		Warnings: migrationWarnings,
	}

	if runtime.GOOS == "linux" {
		installLinuxAssets()
	}

	// Start watching the config directory for changes
	if err := config.WatchConfigDir(ctx); err != nil {
		fmt.Println("Warning: could not start config watcher:", err)
	}
}

// GetApps reads all YAML config files from ~/.config/button/apps/
// and returns them to the frontend. Key resolution is handled client-side
// so the OS toggle can switch without a backend round-trip.
func (a *App) GetApps() (config.AppsResponse, error) {
	return config.ReadApps()
}

// GetCurrentOS returns the detected operating system ("linux" or "darwin").
func (a *App) GetCurrentOS() string {
	return runtime.GOOS
}

// GetAppInfo returns the app name and version from the embedded Wails config.
func (a *App) GetAppInfo() AppInfo {
	type projectInfo struct {
		Name string `json:"name"`
		Info struct {
			ProductName    string `json:"productName"`
			ProductVersion string `json:"productVersion"`
		} `json:"info"`
	}

	result := AppInfo{
		Name: "Button",
	}

	var project projectInfo
	if err := json.Unmarshal(projectConfig, &project); err != nil {
		return result
	}

	if productName := strings.TrimSpace(project.Info.ProductName); productName != "" {
		result.Name = productName
	} else if projectName := strings.TrimSpace(project.Name); projectName != "" {
		result.Name = projectName
	}

	result.Version = strings.TrimSpace(project.Info.ProductVersion)

	return result
}

// GetMigrationResult returns the outcome of the .yml → .yaml migration
// that ran at startup (migrated count + any warnings).
func (a *App) GetMigrationResult() MigrationResult {
	return a.migrationResult
}

// CreateApp creates a new app YAML file in the config directory.
func (a *App) CreateApp(app config.AppConfig) error {
	return config.CreateApp(app)
}

// UpdateApp updates an existing app YAML file. If the name changed and the
// new filename would collide with an existing file, a warning is returned
// (and the write is skipped). Pass force=true to overwrite.
func (a *App) UpdateApp(oldAppName string, app config.AppConfig, force bool) (string, error) {
	return config.UpdateApp(oldAppName, app, force)
}

// DeleteApp removes the YAML file for the given app name.
func (a *App) DeleteApp(appName string) error {
	return config.DeleteApp(appName)
}

// OpenAppFile opens the YAML file for the given app in the system's default text editor.
func (a *App) OpenAppFile(appName string) error {
	path, err := config.AppFilePath(appName)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// GetUserConfig returns the persistent user configuration.
func (a *App) GetUserConfig() config.UserConfig {
	cfg, _ := config.ReadUserConfig()
	return cfg
}

// SetHasSeenWelcome marks the welcome panel as seen in the persistent config.
func (a *App) SetHasSeenWelcome() error {
	cfg, _ := config.ReadUserConfig()
	cfg.HasSeenWelcome = true
	return config.WriteUserConfig(cfg)
}

// SetListPreferences stores persistent app list UI preferences.
func (a *App) SetListPreferences(sortMode string, groupByTag bool) error {
	switch sortMode {
	case "alpha", "last-updated":
	default:
		sortMode = "alpha"
	}

	cfg, _ := config.ReadUserConfig()
	cfg.LastSortMode = sortMode
	cfg.GroupByTag = groupByTag
	return config.WriteUserConfig(cfg)
}

// GetRegistryApps returns the list of apps available in the built-in registry.
func (a *App) GetRegistryApps() ([]registry.RegistryEntry, error) {
	return a.registry.ListApps()
}

// GetExistingAppFiles returns the list of YAML filenames in the user's config directory.
func (a *App) GetExistingAppFiles() ([]string, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	names := make([]string, len(files))
	for i, f := range files {
		names[i] = filepath.Base(f)
	}
	return names, nil
}

// ImportRegistryApps copies selected registry apps to the user's config directory.
// Overwrites existing files if present. Returns the number of apps imported.
func (a *App) ImportRegistryApps(filenames []string) (int, error) {
	dir, err := config.ConfigDir()
	if err != nil {
		return 0, err
	}

	imported := 0
	for _, filename := range filenames {
		destPath := filepath.Join(dir, filename)

		data, err := a.registry.GetAppYAML(filename)
		if err != nil {
			return imported, fmt.Errorf("failed to read registry app %s: %w", filename, err)
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return imported, fmt.Errorf("failed to write %s: %w", filename, err)
		}
		imported++
	}
	return imported, nil
}

// installLinuxAssets writes the app icon and .desktop file to the user's local
// XDG directories so Wayland compositors can resolve the correct taskbar icon.
func installLinuxAssets() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Warning: could not determine home directory:", err)
		return
	}

	iconDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "256x256", "apps")
	if err := os.MkdirAll(iconDir, 0755); err != nil {
		fmt.Println("Warning: could not create icon directory:", err)
		return
	}
	iconPath := filepath.Join(iconDir, "button.png")
	if err := os.WriteFile(iconPath, appIcon, 0644); err != nil {
		fmt.Println("Warning: could not write app icon:", err)
		return
	}

	desktopDir := filepath.Join(home, ".local", "share", "applications")
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		fmt.Println("Warning: could not create applications directory:", err)
		return
	}
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("Warning: could not determine executable path:", err)
		return
	}
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Button
Exec=%s
Icon=button
Categories=Utility;
StartupWMClass=button
`, exe)
	desktopPath := filepath.Join(desktopDir, "button.desktop")
	if err := os.WriteFile(desktopPath, []byte(desktopContent), 0644); err != nil {
		fmt.Println("Warning: could not write .desktop file:", err)
		return
	}

	// Refresh the desktop database so the compositor picks up the changes immediately
	exec.Command("update-desktop-database", desktopDir).Run()
}
