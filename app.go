package main

import (
	"button/internal/config"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Ensure the config directory exists before any reads
	if err := config.EnsureConfigDir(); err != nil {
		fmt.Println("Warning: could not create config directory:", err)
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
