package web

import (
	"fmt"
	"log"
	"time"

	"github.com/jchv/go-webview2"
)

// WebView2App represents our WebView2-based application
type WebView2App struct {
	port       int
	borderInfo *WindowBorderInfo
}

// NewWebView2App creates a new WebView2 application
func NewWebView2App(port int) *WebView2App {
	return &WebView2App{
		port: port,
	}
}

// Run starts the WebView2 application
func (app *WebView2App) Run() error {
	// Start web server in background - WebView2 only, no browser
	go func() {
		server, err := NewServer(app.port)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Create WebView2 instance
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true, // Enable DevTools for layout debugging
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "LaunchRDP",
			Width:  275,
			Height: 750,
			IconId: 2, // Use default icon
		},
	})

	if w == nil {
		return fmt.Errorf("failed to create WebView2")
	}
	defer w.Destroy()

	// Navigate to local server
	url := fmt.Sprintf("http://localhost:%d", app.port)
	w.Navigate(url)

	// Measure window border information after window is shown
	go func() {
		app.measureWindowBorders()
	}()

	fmt.Printf("✓ LaunchRDP WebView2 window opened\n")
	fmt.Printf("Close the window to exit the application\n")

	// Run the WebView2 (this blocks until window is closed)
	w.Run()

	fmt.Printf("WebView2 closed, shutting down...\n")
	return nil
}

// measureWindowBorders measures the window borders after the window is displayed
func (app *WebView2App) measureWindowBorders() {
	// Wait a bit for the window to be fully rendered
	time.Sleep(1 * time.Second)

	// Find our window by title
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		log.Printf("Warning: Could not find LaunchRDP window for border measurement")
		return
	}

	// Measure border information
	borderInfo, err := GetWindowBorderInfo(hwnd)
	if err != nil {
		log.Printf("Warning: Could not measure window borders: %v", err)
		return
	}

	app.borderInfo = borderInfo
	log.Printf("✓ Window border info measured: Left=%d, Top=%d, Right=%d, Bottom=%d, TitleBar=%d",
		borderInfo.Left, borderInfo.Top, borderInfo.Right, borderInfo.Bottom, borderInfo.TitleBarHeight)
}

// GetBorderInfo returns the measured border information
func (app *WebView2App) GetBorderInfo() *WindowBorderInfo {
	return app.borderInfo
}
