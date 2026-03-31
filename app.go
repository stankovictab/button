package main

import (
	"button/internal/config"
	"context"
	"fmt"
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
