package web

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/jchv/go-webview2"
)

// getWorkArea returns the work area (usable desktop area) coordinates
func getWorkArea() (x, y, width, height int32, err error) {
	// Use SystemParametersInfo to get work area (desktop minus taskbar)
	procSystemParametersInfo := user32.NewProc("SystemParametersInfoW")

	var rect RECT
	ret, _, err := procSystemParametersInfo.Call(
		uintptr(48), // SPI_GETWORKAREA
		uintptr(0),
		uintptr(unsafe.Pointer(&rect)),
		uintptr(0),
	)

	if ret == 0 {
		return 0, 0, 0, 0, fmt.Errorf("failed to get work area: %v", err)
	}

	return rect.Left, rect.Top, rect.Right - rect.Left, rect.Bottom - rect.Top, nil
}

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

	// Get work area dimensions (usable desktop area, excluding taskbar)
	workX, workY, _, workHeight, err := getWorkArea()
	if err != nil {
		log.Printf("Failed to get work area, using defaults: %v", err)
		workX, workY, workHeight = 0, 0, 750
	}

	// Create WebView2 instance with full work area height
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true, // Enable DevTools for layout debugging
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "LaunchRDP",
			Width:  275,              // Keep fixed width for usability
			Height: uint(workHeight), // Use full work area height
			IconId: 2,                // Use default icon
		},
	})

	if w == nil {
		return fmt.Errorf("failed to create WebView2")
	}
	defer w.Destroy()

	// Position window at left edge of work area
	app.positionWindow(w, workX, workY)

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

// positionWindow positions the WebView2 window at the specified coordinates
func (app *WebView2App) positionWindow(w webview2.WebView, x, y int32) {
	// Wait a moment for the window to be created
	time.Sleep(100 * time.Millisecond)

	// Find our window by title
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		log.Printf("Warning: Could not find LaunchRDP window for positioning")
		return
	}

	// Use SetWindowPos to position the window
	procSetWindowPos := user32.NewProc("SetWindowPos")
	ret, _, err := procSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(0),   // HWND_TOP
		uintptr(x),   // X position
		uintptr(y),   // Y position
		uintptr(0),   // Width (0 = no change)
		uintptr(0),   // Height (0 = no change)
		uintptr(0x1), // SWP_NOSIZE (don't change size)
	)

	if ret == 0 {
		log.Printf("Warning: Failed to position window: %v", err)
	} else {
		log.Printf("✓ Window positioned at (%d, %d)", x, y)
	}
}
