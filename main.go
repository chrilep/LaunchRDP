package main

import (
	"flag"

	"github.com/chrilep/LaunchRDP/config"
	"github.com/chrilep/LaunchRDP/logging"
	"github.com/chrilep/LaunchRDP/web"
)

func main() {
	debug := false

	// Set up panic handling and logging
	defer logging.PanicHandler()

	// Initialize logging
	logging.Log(debug, "===", AppName, "Starting ===")
	logging.Log(debug, "Version:", Version)
	logging.Log(debug, "Build Date: 2025-10-29")

	// Command line flags
	var port = flag.Int("port", 9457, "Port for web server") // Use consistent port to avoid firewall prompts
	var version = flag.Bool("version", false, "Show version information")
	var v = flag.Bool("v", false, "Show version information")
	flag.Parse()

	// Check for version flag
	if *version || *v {
		PrintVersion()
		return
	}

	logging.Log(debug, "Initializing web-based LaunchRDP")

	// Initialize config directories
	logging.Log(debug, "Initializing config directories")
	if err := config.InitDirectories(); err != nil {
		logging.Log(true, "FATAL: Failed to initialize directories:", err)
		return // Exit gracefully instead of fatal error
	}
	logging.Log(debug, "Config directories initialized successfully")
	logging.Log(debug, "Starting WebView2-only application")

	// Create and run WebView2 app (WebView2-only, no browser support)
	app := web.NewWebView2App(*port)
	if err := app.Run(); err != nil {
		logging.Log(true, "FATAL: WebView2 failed:", err)
		return // Exit gracefully instead of fatal error
	}

	// This should only be reached when the application exits
	logging.Log(debug, "===", AppName, "Exiting ===")
	logging.CloseLog()
}
