package main

import (
	"embed"
	"flag"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	// Your existing packages for initialization
	"github.com/chrilep/LaunchRDP/app/config"
	"github.com/chrilep/LaunchRDP/app/logging"
)

//go:embed all:frontend/dist
var assets embed.FS

// Icon can be added later when build/appicon.png exists
// var icon []byte

func main() {
	debug := true

	// Set up panic handling and logging - SAME AS BEFORE!
	defer logging.PanicHandler()

	// Initialize logging - SAME AS BEFORE!
	logging.Log(debug, "===", AppName, "Starting (Wails Version) ===")
	logging.Log(debug, "Version:", Version)
	logging.Log(debug, "Build Date: 2025-11-07")

	// Command line flags - SAME AS BEFORE!
	var version = flag.Bool("version", false, "Show version information")
	var v = flag.Bool("v", false, "Show version information")
	flag.Parse()

	// Check for version flag - SAME AS BEFORE!
	if *version || *v {
		PrintVersion()
		return
	}

	logging.Log(debug, "Initializing Wails-based LaunchRDP")

	// Initialize config directories - SAME AS BEFORE!
	logging.Log(debug, "Initializing config directories")
	if err := config.InitDirectories(); err != nil {
		logging.Log(true, "FATAL: Failed to initialize directories:", err)
		return // Exit gracefully instead of fatal error
	}
	logging.Log(debug, "Config directories initialized successfully")
	logging.Log(debug, "Starting Wails application")

	// Create an instance of the app structure
	app := NewLaunchRDPApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "LaunchRDP",
		Width:  275,
		Height: 1379, // 1349 + 30 for MyDockFinder

		// Window positioning
		MinWidth:           275,
		MinHeight:          500,
		StartHidden:        true, // Anti-flicker
		HideWindowOnClose:  false,
		DisableResize:      false,
		Fullscreen:         false,
		Frameless:          false,
		SingleInstanceLock: &options.SingleInstanceLock{UniqueId: ID},

		// Styling
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},

		// Assets
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		// Lifecycle hooks
		OnStartup:  app.Startup,
		OnDomReady: app.DomReady,
		OnShutdown: app.Shutdown,

		// Enable development features in debug mode
		Debug: options.Debug{
			OpenInspectorOnStartup: debug,
		},

		// Bind exported Go structures/methods to the frontend so wailsjs/go is generated
		Bind: []interface{}{app},
	})

	if err != nil {
		logging.Log(true, "FATAL: Wails application failed:", err)
		return // Exit gracefully instead of fatal error
	}
}
