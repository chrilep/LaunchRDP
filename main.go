package main

import (
	"flag"
	"log"

	"github.com/chrilep/LaunchRDP/config"
	"github.com/chrilep/LaunchRDP/logging"
	"github.com/chrilep/LaunchRDP/web"
)

func main() {
	// Set up panic handling and logging
	defer logging.PanicHandler()

	// Initialize logging
	logging.Log(true, "===", AppName, "Starting ===")
	logging.Log(true, "Version:", Version)
	logging.Log(true, "Build Date: 2025-10-29")

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

	logging.Log(true, "Initializing web-based LaunchRDP")

	// Initialize config directories
	logging.Log(true, "Initializing config directories")
	if err := config.InitDirectories(); err != nil {
		logging.Log(true, "ERROR: Failed to initialize directories:", err)
		log.Fatalf("Failed to initialize directories: %v", err)
	}
	logging.Log(true, "Config directories initialized successfully")
	logging.Log(true, "Starting WebView2-only application")

	// Create and run WebView2 app (WebView2-only, no browser support)
	app := web.NewWebView2App(*port)
	if err := app.Run(); err != nil {
		logging.Log(true, "ERROR: WebView2 failed:", err)
		log.Fatalf("Failed to start WebView2: %v", err)
	}

	// This should only be reached when the application exits
	logging.Log(true, "===", AppName, "Exiting ===")
	logging.CloseLog()
}
