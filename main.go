package main

import (
	"button/internal/config"
	"button/internal/linuxtray"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

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

//go:embed assets/images/button-logo-systray.png
var trayIcon []byte

//go:embed wails.json
var projectConfig []byte

//go:embed registry/*.yaml
var registryFS embed.FS

const singleInstanceID = "button.stankovictab.button"

func main() {
	action, err := parseLaunchAction(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(2)
	}
	if handled := handleExistingLinuxInstance(singleInstanceID, action, os.Args[1:]); handled {
		return
	}

	// The registry/ embed includes the "registry" prefix in paths.
	// Sub into it so the embedded registry sees flat filenames.
	regFS, err := fs.Sub(registryFS, "registry")
	if err != nil {
		log.Fatal("failed to load registry:", err)
	}
	builtInRegistry := config.NewEmbeddedRegistry(regFS)

	// Create an instance of the app structure
	app := NewApp(builtInRegistry, action, linuxtray.New(trayIcon))

	// Create application with options
	err = wails.Run(&options.App{
		Title:       "Button",
		Width:       900,
		Height:      600,
		StartHidden: action == launchQuit,
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
		OnDomReady:       app.domReady,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    app.beforeClose,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: singleInstanceID,
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				app.handleSecondInstanceLaunch(secondInstanceData.Args)
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

type launchAction string

const (
	launchDefault launchAction = "default"
	launchToggle  launchAction = "toggle"
	launchQuit    launchAction = "quit"
)

func parseLaunchAction(args []string) (launchAction, error) {
	switch len(args) {
	case 0:
		return launchDefault, nil
	case 1:
		switch args[0] {
		case "--toggle":
			return launchToggle, nil
		case "--quit":
			return launchQuit, nil
		}
	}

	return "", fmt.Errorf("invalid arguments")
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: button [--toggle|--quit]")
}
