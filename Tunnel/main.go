package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	app := NewApp()

	tempDir, err := os.MkdirTemp("", "wails_userdata_*")
	if err != nil {
		panic("Failed to create temp user data directory: " + err.Error())
	}

	defer os.RemoveAll(tempDir)

	// Create application with options
	err = wails.Run(&options.App{
		Title:    "Tunnel",
		Width:    1200,
		Height:   768,
		MinWidth: 1100,
		// Frameless: true,
		Windows: &windows.Options{
			Theme: windows.Dark,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar: windows.RGB(255, 255, 255),
			},
			IsZoomControlEnabled: false,
			WebviewUserDataPath:  tempDir,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},

		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
