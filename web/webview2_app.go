package web

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/chrilep/LaunchRDP/logging"
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
	port        int
	borderInfo  *WindowBorderInfo
	minWidth    int
	maxWidth    int
	stopMonitor chan bool
}

// NewWebView2App creates a new WebView2 application
func NewWebView2App(port int) *WebView2App {
	return &WebView2App{
		port:        port,
		stopMonitor: make(chan bool),
	}
}

// Run starts the WebView2 application
func (app *WebView2App) Run() error {
	debug := false

	// Start web server in background - WebView2 only, no browser
	go func() {
		server, err := NewServer(app.port)
		if err != nil {
			logging.Log(true, "FATAL: Failed to create server:", err)
			return
		}
		if err := server.Start(); err != nil {
			logging.Log(true, "Server error:", err)
		}
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Get work area dimensions (usable desktop area, excluding taskbar)
	workX, workY, _, workHeight, err := getWorkArea()
	if err != nil {
		logging.Log(true, "Failed to get work area, using defaults:", err)
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
	app.positionWindow(workX, workY)

	// Set window size constraints (min: 275px, max: 500px width)
	app.setWindowSizeConstraints(275, 500)

	// Navigate to local server
	url := fmt.Sprintf("http://localhost:%d", app.port)
	w.Navigate(url)

	// Measure window border information after window is shown
	go func() {
		app.measureWindowBorders()
	}()

	logging.Log(debug, "✓ LaunchRDP WebView2 window opened")
	logging.Log(debug, "Close the window to exit the application")

	// Run the WebView2 (this blocks until window is closed)
	w.Run()

	// Stop the window size monitor
	close(app.stopMonitor)

	logging.Log(debug, "WebView2 closed, shutting down...")
	return nil
}

// measureWindowBorders measures the window borders after the window is displayed
func (app *WebView2App) measureWindowBorders() {
	debug := false

	// Wait a bit for the window to be fully rendered
	time.Sleep(1 * time.Second)

	// Find our window by title
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		logging.Log(true, "Warning: Could not find LaunchRDP window for border measurement")
		return
	}

	// Measure border information
	borderInfo, err := GetWindowBorderInfo(hwnd)
	if err != nil {
		logging.Log(true, "Warning: Could not measure window borders:", err)
		return
	}

	app.borderInfo = borderInfo
	logging.Log(debug, "✓ Window border info measured: Left=", borderInfo.Left, "Top=", borderInfo.Top, "Right=", borderInfo.Right, "Bottom=", borderInfo.Bottom, "TitleBar=", borderInfo.TitleBarHeight)
}

// GetBorderInfo returns the measured border information
func (app *WebView2App) GetBorderInfo() *WindowBorderInfo {
	return app.borderInfo
}

// positionWindow positions the WebView2 window at the specified coordinates
func (app *WebView2App) positionWindow(x, y int32) {
	debug := false

	// Wait a moment for the window to be created
	time.Sleep(100 * time.Millisecond)

	// Find our window by title
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		logging.Log(true, "Warning: Could not find LaunchRDP window for positioning")
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
		logging.Log(true, "Warning: Failed to position window:", err)
		return
	}

	logging.Log(debug, "✓ Window positioned at (", x, ",", y, ")")
}

// setWindowSizeConstraints sets minimum and maximum width constraints for the window
func (app *WebView2App) setWindowSizeConstraints(minWidth, maxWidth int) {
	debug := false

	// Wait a moment for the window to be created
	time.Sleep(100 * time.Millisecond)

	// Find our window by title
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		logging.Log(true, "Warning: Could not find LaunchRDP window for size constraints")
		return
	}

	// We need to subclass the window to handle WM_GETMINMAXINFO messages
	// This is more complex, so let's use a simpler approach with SetWindowLongPtr
	// to store our constraints and handle them in a window proc

	// For now, let's just log that constraints are being set
	// The actual implementation would require setting up a custom window procedure
	logging.Log(debug, "✓ Window size constraints set: min=", minWidth, "px, max=", maxWidth, "px width")

	// Store the constraints in the app
	app.minWidth = minWidth
	app.maxWidth = maxWidth

	// Start monitoring window size in a goroutine
	go app.monitorWindowSize()
}

// monitorWindowSize monitors the window size and enforces constraints
func (app *WebView2App) monitorWindowSize() {
	debug := false

	ticker := time.NewTicker(50 * time.Millisecond) // Faster polling for better responsiveness
	defer ticker.Stop()

	logging.Log(debug, "✓ Started window size monitoring with constraints:", app.minWidth, "-", app.maxWidth, "px width")

	for {
		select {
		case <-ticker.C:
			app.enforceWindowConstraints()
		case <-app.stopMonitor:
			logging.Log(debug, "✓ Stopped window size monitoring")
			return
		}
	}
}

// enforceWindowConstraints checks and corrects window size if needed
func (app *WebView2App) enforceWindowConstraints() {
	debug := false
	hwnd := findWindowByTitle("LaunchRDP")
	if hwnd == 0 {
		logging.Log(true, "ERROR: Cannot find my window with title LaunchRDP")
		return
	}
	logging.Log(debug, "Found window handle:", hwnd)

	// Get current window rectangle
	var rect RECT
	procGetWindowRect := user32.NewProc("GetWindowRect")
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		logging.Log(true, "ERROR: Failed to get window rect")
		return
	}
	logging.Log(debug, "Current window rect:", rect)

	currentWidth := int(rect.Right - rect.Left)
	currentHeight := int(rect.Bottom - rect.Top)
	var newWidth int
	needsResize := false

	// Check if width violates constraints
	if currentWidth < app.minWidth {
		newWidth = app.minWidth
		needsResize = true
		logging.Log(true, "Window too narrow (", currentWidth, "px), resizing to minimum:", newWidth, "px")
	} else if currentWidth > app.maxWidth {
		newWidth = app.maxWidth
		needsResize = true
		logging.Log(true, "Window too wide (", currentWidth, "px), resizing to maximum:", newWidth, "px")
	}

	// Resize if needed
	if needsResize {
		procSetWindowPos := user32.NewProc("SetWindowPos")
		ret, _, err := procSetWindowPos.Call(
			uintptr(hwnd),
			uintptr(0),             // HWND_TOP
			uintptr(rect.Left),     // X position (keep current)
			uintptr(rect.Top),      // Y position (keep current)
			uintptr(newWidth),      // New width
			uintptr(currentHeight), // Height (keep current)
			uintptr(0x0044),        // SWP_NOZORDER | SWP_SHOWWINDOW
		)
		if ret == 0 {
			logging.Log(true, "ERROR: Failed to resize window:", err)
			return
		}
		logging.Log(debug, "Window resized successfully to", newWidth, "x", currentHeight)
	}
}
