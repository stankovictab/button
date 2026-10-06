package main

import (
	"button/internal/config"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBeforeCloseRequiresEnabledWorkingTray(t *testing.T) {
	for _, tc := range []struct {
		name                              string
		enabled, ready, quit, wantPrevent bool
	}{
		{"disabled", false, true, false, false},
		{"tray failed", true, false, false, false},
		{"close to tray", true, true, false, runtime.GOOS == "linux"},
		{"explicit quit", true, true, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := &App{closeToTray: tc.enabled, trayReady: tc.ready, allowQuit: tc.quit}
			if got := a.beforeClose(context.Background()); got != tc.wantPrevent {
				t.Fatalf("prevent close = %v, want %v", got, tc.wantPrevent)
			}
		})
	}
}

func TestCloseToTrayPersistsAndPreservesOtherPreferences(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux-only preference")
	}
	t.Setenv("HOME", t.TempDir())
	if err := config.WriteUserConfig(config.UserConfig{HasSeenWelcome: true, LastSortMode: "name", GroupByTag: true}); err != nil {
		t.Fatal(err)
	}
	a := &App{}
	for _, enabled := range []bool{true, false} {
		if err := a.SetCloseToTray(enabled); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.ReadUserConfig()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.CloseToTray != enabled || a.closeToTray != enabled || !cfg.HasSeenWelcome || cfg.LastSortMode != "name" || !cfg.GroupByTag {
			t.Fatalf("preference update lost state: %+v", cfg)
		}
	}
}

func TestCloseToTrayReadFailureDoesNotEnableOrOverwrite(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux-only preference")
	}
	t.Setenv("HOME", t.TempDir())
	dir, _ := config.BaseDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("closeToTray: ["), 0644); err != nil {
		t.Fatal(err)
	}
	a := &App{}
	if err := a.SetCloseToTray(true); err == nil {
		t.Fatal("expected invalid config error")
	}
	if a.closeToTray {
		t.Fatal("enabled preference after failed write")
	}
}
