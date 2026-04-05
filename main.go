package main

import (
	"button/internal/registry"
	"embed"
	"io/fs"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

//go:embed wails.json
var projectConfig []byte

//go:embed registry/*.yaml
var registryFS embed.FS

func main() {
	// The registry/ embed includes the "registry" prefix in paths.
	// Sub into it so the Registry sees flat filenames.
	regFS, err := fs.Sub(registryFS, "registry")
	if err != nil {
		log.Fatal("failed to load registry:", err)
	}
	reg := registry.New(regFS)

	// Create an instance of the app structure
	app := NewApp(reg)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Button",
		Width:  900,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Linux: &linux.Options{
			Icon:        appIcon,
			ProgramName: "button",
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		BackgroundColour: &options.RGBA{R: 17, G: 17, B: 17, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
