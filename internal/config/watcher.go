package config

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// WatchConfigDir watches ~/.config/button/apps/ for YAML file changes
// and emits a "config:changed" Wails event whenever a .yaml file is
// created, modified, removed, or renamed. Uses a short debounce to
// coalesce rapid successive writes (e.g. editor save-then-rename).
func WatchConfigDir(ctx context.Context) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch %s: %w", dir, err)
	}

	go func() {
		defer watcher.Close()

		var debounce *time.Timer

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Only react to .yaml/.yml files
				name := filepath.Base(event.Name)
				if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
					continue
				}
				// Debounce: wait 200ms for rapid successive events to settle
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(200*time.Millisecond, func() {
					resp, err := ReadApps()
					if err != nil {
						fmt.Println("config watcher: error reading apps:", err)
						return
					}
					wailsRuntime.EventsEmit(ctx, "config:changed", resp)
				})
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("config watcher error:", err)
			}
		}
	}()

	return nil
}
