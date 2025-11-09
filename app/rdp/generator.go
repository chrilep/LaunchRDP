package rdp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/chrilep/LaunchRDP/app/config"
	"github.com/chrilep/LaunchRDP/app/logging"
	"github.com/chrilep/LaunchRDP/app/models"
)

// Windows API declarations for window enumeration
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procGetClassNameW       = user32.NewProc("GetClassNameW")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procShowWindow          = user32.NewProc("ShowWindow")
	procIsIconic            = user32.NewProc("IsIconic")
)

const (
	SW_RESTORE = 9
)

// Generator handles RDP file generation and launching
type Generator struct {
	// Callback function to save user after password migration
	SaveUserCallback func(user models.User) error
}

// NewGenerator creates a new RDP generator
func NewGenerator() *Generator {
	return &Generator{}
}

// SetSaveUserCallback sets the callback function for saving users after password migration
func (g *Generator) SetSaveUserCallback(callback func(user models.User) error) {
	g.SaveUserCallback = callback
}

// GenerateRDPFile creates a temporary RDP file with the specified settings
func (g *Generator) GenerateRDPFile(host models.Host, user models.User) (string, error) {
	debug := false
	logging.Log(debug, "GenerateRDPFile started for host:", host.Name, "user:", user.Username)

	// Use host ID for safe filename - avoids issues with special characters in usernames/hostnames
	filename := fmt.Sprintf("%s.rdp", host.ID)
	filepath := config.GetTempPath(filename)
	logging.Log(debug, "RDP file path (using host ID for safety):", filepath)

	// Build RDP file content
	logging.Log(debug, "Building RDP content")
	content := g.buildRDPContent(host, user)
	logging.Log(debug, "RDP content built, length:", len(content), "bytes")

	// Write to temporary file
	logging.Log(debug, "Writing RDP file to disk")
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		logging.Log(true, "ERROR: Failed to write RDP file:", err)
		return "", fmt.Errorf("failed to write RDP file: %w", err)
	}
	logging.Log(debug, "RDP file written successfully")

	return filepath, nil
}

// buildRDPContent creates the RDP file content based on host and user settings
func (g *Generator) buildRDPContent(host models.Host, user models.User) string {
	debug := false
	var builder strings.Builder

	// Core connection settings
	builder.WriteString(fmt.Sprintf("full address:s:%s\n", host.Address))
	builder.WriteString(fmt.Sprintf("server port:i:%d\n", host.Port))

	// Username
	builder.WriteString(fmt.Sprintf("username:s:%s\n", user.Username))

	// Redirection & display base
	if host.RedirectClipboard {
		builder.WriteString("redirectclipboard:i:1\n")
	} else {
		builder.WriteString("redirectclipboard:i:0\n")
	}
	// Dynamic resolution off for fixed window; could be toggled later
	builder.WriteString("dynamic resolution:i:0\n")
	// Screen mode: 2 = fullscreen, 1 = windowed
	modeID := 1
	if strings.ToLower(host.DisplayMode) == "fullscreen" || host.ScreenMode == 2 {
		modeID = 2
	}
	builder.WriteString(fmt.Sprintf("screen mode id:i:%d\n", modeID))

	if modeID == 2 {
		// Fullscreen: use current desktop metrics approximated by large size if none stored
		w := host.DesktopWidth
		h := host.DesktopHeight
		if w <= 0 || h <= 0 {
			w = 1920
			h = 1080
		}
		builder.WriteString(fmt.Sprintf("desktopwidth:i:%d\n", w))
		builder.WriteString(fmt.Sprintf("desktopheight:i:%d\n", h))
	} else {
		builder.WriteString(fmt.Sprintf("desktopwidth:i:%d\n", host.DesktopWidth))
		builder.WriteString(fmt.Sprintf("desktopheight:i:%d\n", host.DesktopHeight))
	}

	// Window positioning - calculate winposstr from current values
	// Format: "0,1,<x>,<y>,<right>,<bottom>"
	windowWidth := host.WindowWidth
	windowHeight := host.WindowHeight

	// DEBUG: Log all host values
	logging.Log(debug, "RDP Generator Debug - Host Values:")
	logging.Log(debug, "  WindowWidth:", windowWidth)
	logging.Log(debug, "  WindowHeight:", windowHeight)
	logging.Log(debug, "  DesktopWidth:", host.DesktopWidth)
	logging.Log(debug, "  DesktopHeight:", host.DesktopHeight)
	logging.Log(debug, "  PositionX:", host.PositionX)
	logging.Log(debug, "  PositionY:", host.PositionY)

	// Fallback if window dimensions are not set (legacy data or missing values)
	if windowWidth == 0 || windowHeight == 0 {
		logging.Log(debug, "  Using fallback calculation for window size")
		// Estimate window size from desktop size (add typical border sizes)
		windowWidth = host.DesktopWidth + 16   // Add estimated horizontal borders
		windowHeight = host.DesktopHeight + 59 // Add estimated title bar + borders
		logging.Log(debug, "  Calculated WindowWidth:", windowWidth)
		logging.Log(debug, "  Calculated WindowHeight:", windowHeight)
	}

	windowRight := host.PositionX + windowWidth
	windowBottom := host.PositionY + windowHeight
	winPosStr := fmt.Sprintf("0,1,%d,%d,%d,%d", host.PositionX, host.PositionY, windowRight, windowBottom)

	logging.Log(debug, "  Final winPosStr:", winPosStr)
	builder.WriteString(fmt.Sprintf("winposstr:s:%s\n", winPosStr))

	// Display and performance settings
	// Multi-monitor support: enable when fullscreen, disable for windowed mode
	if modeID == 2 {
		builder.WriteString("use multimon:i:1\n")
	} else {
		builder.WriteString("use multimon:i:0\n")
	}
	builder.WriteString("session bpp:i:32\n")
	builder.WriteString("compression:i:1\n")
	builder.WriteString("keyboardhook:i:1\n")
	builder.WriteString("audiocapturemode:i:1\n")
	builder.WriteString("videoplaybackmode:i:1\n")
	builder.WriteString("connection type:i:7\n")
	builder.WriteString("networkautodetect:i:1\n")
	builder.WriteString("bandwidthautodetect:i:1\n")
	builder.WriteString("displayconnectionbar:i:1\n")
	builder.WriteString("enableworkspacereconnect:i:0\n")
	builder.WriteString("remoteappmousemoveinject:i:1\n")

	// Visual performance settings
	builder.WriteString("disable wallpaper:i:0\n")
	builder.WriteString("allow font smoothing:i:0\n")
	builder.WriteString("allow desktop composition:i:0\n")
	builder.WriteString("disable full window drag:i:1\n")
	builder.WriteString("disable menu anims:i:1\n")
	builder.WriteString("disable themes:i:0\n")
	builder.WriteString("disable cursor setting:i:0\n")
	builder.WriteString("bitmapcachepersistenable:i:1\n")

	// Audio and device redirection
	builder.WriteString("audiomode:i:0\n")
	builder.WriteString("redirectprinters:i:1\n")
	builder.WriteString("redirectlocation:i:1\n")
	builder.WriteString("redirectcomports:i:1\n")
	builder.WriteString("redirectsmartcards:i:1\n")
	builder.WriteString("redirectwebauthn:i:1\n")
	builder.WriteString("redirectposdevices:i:0\n")
	builder.WriteString("camerastoredirect:s:*\n")
	builder.WriteString("devicestoredirect:s:*\n")

	// Drive redirection (configurable)
	if host.RedirectDrives {
		builder.WriteString("drivestoredirect:s:*\n")
	} else {
		builder.WriteString("drivestoredirect:s:\n")
	}

	// Connection and security settings
	builder.WriteString("autoreconnection enabled:i:1\n")
	builder.WriteString("authentication level:i:2\n")
	builder.WriteString("prompt for credentials:i:0\n")
	builder.WriteString("negotiate security layer:i:1\n")
	builder.WriteString("remoteapplicationmode:i:0\n")
	builder.WriteString("alternate shell:s:\n")
	builder.WriteString("shell working directory:s:\n")

	// Gateway settings (empty by default)
	builder.WriteString("gatewayhostname:s:\n")
	builder.WriteString("gatewayusagemethod:i:4\n")
	builder.WriteString("gatewaycredentialssource:i:4\n")
	builder.WriteString("gatewayprofileusagemethod:i:0\n")
	builder.WriteString("promptcredentialonce:i:0\n")
	builder.WriteString("gatewaybrokeringtype:i:0\n")
	builder.WriteString("use redirection server name:i:0\n")
	builder.WriteString("rdgiskdcproxy:i:0\n")
	builder.WriteString("kdcproxyname:s:\n")
	builder.WriteString("enablerdsaadauth:i:0\n")

	return builder.String()
}

// findExistingRDPWindow searches for an existing mstsc.exe window with the target address
func findExistingRDPWindow(targetAddress string) (uintptr, bool) {
	debug := false
	logging.Log(debug, "Searching for existing RDP window for:", targetAddress)

	var foundHwnd uintptr

	// Callback function for EnumWindows
	callback := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		// Check if window is visible
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1 // Continue enumeration
		}

		// Get window class name
		className := make([]uint16, 256)
		procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&className[0])), 256)
		classNameStr := syscall.UTF16ToString(className)

		// Check if it's a Terminal Services Client window (mstsc.exe)
		if classNameStr != "TscShellContainerClass" {
			return 1 // Continue enumeration
		}

		// Get window title
		titleBuf := make([]uint16, 512)
		procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), 512)
		title := syscall.UTF16ToString(titleBuf)

		logging.Log(debug, "Found RDP window with title:", title)

		// Check if the title contains the target address
		if strings.Contains(strings.ToLower(title), strings.ToLower(targetAddress)) {
			logging.Log(debug, "Match found! HWND:", hwnd)
			foundHwnd = hwnd
			return 0 // Stop enumeration
		}

		return 1 // Continue enumeration
	})

	// Enumerate all windows
	procEnumWindows.Call(callback, 0)

	if foundHwnd != 0 {
		logging.Log(debug, "Existing RDP window found:", foundHwnd)
		return foundHwnd, true
	}

	logging.Log(debug, "No existing RDP window found")
	return 0, false
}

// bringWindowToFront brings a window to the foreground and restores it if minimized
func bringWindowToFront(hwnd uintptr) {
	debug := false
	logging.Log(debug, "Bringing window to front, HWND:", hwnd)

	// Check if window is minimized
	isMinimized, _, _ := procIsIconic.Call(hwnd)
	if isMinimized != 0 {
		logging.Log(debug, "Window is minimized, restoring...")
		procShowWindow.Call(hwnd, SW_RESTORE)
	}

	// Bring window to foreground
	procSetForegroundWindow.Call(hwnd)
	logging.Log(debug, "Window brought to front")
}

// LaunchRDP launches an RDP session using mstsc.exe
func (g *Generator) LaunchRDP(rdpFilePath string) error {
	debug := false
	logging.Log(debug, "LaunchRDP started with file:", rdpFilePath)

	// Use absolute path
	logging.Log(debug, "Getting absolute path for RDP file")
	absPath, err := filepath.Abs(rdpFilePath)
	if err != nil {
		logging.Log(true, "ERROR: Failed to get absolute path:", err)
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	logging.Log(debug, "Absolute path:", absPath)

	// Launch mstsc.exe with the RDP file
	logging.Log(debug, "Launching mstsc.exe with RDP file")
	cmd := exec.Command("mstsc", absPath)

	logging.Log(debug, "Starting mstsc process...")
	if err := cmd.Start(); err != nil {
		logging.Log(true, "ERROR: Failed to launch mstsc:", err)
		return fmt.Errorf("failed to launch mstsc: %w", err)
	}

	logging.Log(debug, "mstsc process started successfully, PID:", cmd.Process.Pid)
	return nil
}

// LaunchHost launches an RDP session for the specified host and user
// Returns (wasReused bool, error) - wasReused is true if existing window was activated
func (g *Generator) LaunchHost(host models.Host, user models.User) (bool, error) {
	debug := false
	logging.Log(debug, "LaunchHost started for host:", host.Name, "address:", host.Address, "user:", user.Username)
	logging.Log(debug, "User details - ID:", user.ID, "Name:", user.Name, "Username:", user.Username)
	logging.Log(debug, "User has encrypted password:", user.EncryptedPassword != "")

	// Check for existing RDP window first
	if hwnd, found := findExistingRDPWindow(host.Address); found {
		logging.Log(true, "Found existing RDP window for", host.Address, "- bringing to front")
		bringWindowToFront(hwnd)
		return true, nil
	}

	// User password should already be stored in Windows CredStore when user was saved
	if user.EncryptedPassword != "" {
		logging.Log(debug, "User has encrypted password stored")
	}

	logging.Log(debug, "Using existing credentials from Windows Credential Store (stored during user save)") // Generate standard RDP file (password comes from CredStore)
	logging.Log(debug, "Generating RDP file")
	rdpFile, err := g.GenerateRDPFile(host, user)
	if err != nil {
		logging.Log(true, "ERROR: Failed to generate RDP file:", err)
		return false, fmt.Errorf("failed to generate RDP file: %w", err)
	}
	logging.Log(debug, "RDP file generated:", rdpFile)

	// Launch RDP session
	logging.Log(debug, "Launching RDP session")
	if err := g.LaunchRDP(rdpFile); err != nil {
		logging.Log(true, "ERROR: Failed to launch RDP session:", err)
		return false, fmt.Errorf("failed to launch RDP session: %w", err)
	}

	logging.Log(debug, "LaunchHost completed successfully")
	return false, nil
}

// CleanupTempFiles removes old RDP files from temp directory
func (g *Generator) CleanupTempFiles() error {
	tempDir := config.TempDir

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rdp") {
			filePath := filepath.Join(tempDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				// Log error but continue cleanup
				logging.Log(true, "Warning: failed to remove temp file", filePath+":", err)
			}
		}
	}

	return nil
}
