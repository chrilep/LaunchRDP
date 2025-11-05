package rdp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chrilep/LaunchRDP/config"
	"github.com/chrilep/LaunchRDP/logging"
	"github.com/chrilep/LaunchRDP/models"
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
	logging.Log(true, "GenerateRDPFile started for host:", host.Name, "user:", user.Username)

	// Use host ID for safe filename - avoids issues with special characters in usernames/hostnames
	filename := fmt.Sprintf("%s.rdp", host.ID)
	filepath := config.GetTempPath(filename)
	logging.Log(true, "RDP file path (using host ID for safety):", filepath)

	// Build RDP file content
	logging.Log(true, "Building RDP content")
	content := g.buildRDPContent(host, user)
	logging.Log(true, "RDP content built, length:", len(content), "bytes")

	// Write to temporary file
	logging.Log(true, "Writing RDP file to disk")
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		logging.Log(true, "ERROR: Failed to write RDP file:", err)
		return "", fmt.Errorf("failed to write RDP file: %w", err)
	}
	logging.Log(true, "RDP file written successfully")

	return filepath, nil
}

// buildRDPContent creates the RDP file content based on host and user settings
func (g *Generator) buildRDPContent(host models.Host, user models.User) string {
	var builder strings.Builder

	// Core connection settings
	builder.WriteString(fmt.Sprintf("full address:s:%s\n", host.Address))
	builder.WriteString(fmt.Sprintf("server port:i:%d\n", host.Port))
	builder.WriteString(fmt.Sprintf("username:s:%s\n", user.Username))

	// Standard RDP settings - comprehensive template
	builder.WriteString("redirectclipboard:i:1\n")
	builder.WriteString("dynamic resolution:i:0\n")
	builder.WriteString("screen mode id:i:1\n")

	// Desktop dimensions (calculated from window size minus borders)
	builder.WriteString(fmt.Sprintf("desktopwidth:i:%d\n", host.DesktopWidth))
	builder.WriteString(fmt.Sprintf("desktopheight:i:%d\n", host.DesktopHeight))

	// Window positioning - calculate winposstr from current values
	// Format: "0,1,<x>,<y>,<right>,<bottom>"
	windowWidth := host.WindowWidth
	windowHeight := host.WindowHeight

	// DEBUG: Log all host values
	logging.Log(true, "RDP Generator Debug - Host Values:")
	logging.Log(true, "  WindowWidth:", windowWidth)
	logging.Log(true, "  WindowHeight:", windowHeight)
	logging.Log(true, "  DesktopWidth:", host.DesktopWidth)
	logging.Log(true, "  DesktopHeight:", host.DesktopHeight)
	logging.Log(true, "  PositionX:", host.PositionX)
	logging.Log(true, "  PositionY:", host.PositionY)

	// Fallback if window dimensions are not set (legacy data or missing values)
	if windowWidth == 0 || windowHeight == 0 {
		logging.Log(true, "  Using fallback calculation for window size")
		// Estimate window size from desktop size (add typical border sizes)
		windowWidth = host.DesktopWidth + 16   // Add estimated horizontal borders
		windowHeight = host.DesktopHeight + 59 // Add estimated title bar + borders
		logging.Log(true, "  Calculated WindowWidth:", windowWidth)
		logging.Log(true, "  Calculated WindowHeight:", windowHeight)
	}

	windowRight := host.PositionX + windowWidth
	windowBottom := host.PositionY + windowHeight
	winPosStr := fmt.Sprintf("0,1,%d,%d,%d,%d", host.PositionX, host.PositionY, windowRight, windowBottom)

	logging.Log(true, "  Final winPosStr:", winPosStr)
	builder.WriteString(fmt.Sprintf("winposstr:s:%s\n", winPosStr))

	// Display and performance settings
	builder.WriteString("use multimon:i:0\n")
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

// LaunchRDP launches an RDP session using mstsc.exe
func (g *Generator) LaunchRDP(rdpFilePath string) error {
	logging.Log(true, "LaunchRDP started with file:", rdpFilePath)

	// Use absolute path
	logging.Log(true, "Getting absolute path for RDP file")
	absPath, err := filepath.Abs(rdpFilePath)
	if err != nil {
		logging.Log(true, "ERROR: Failed to get absolute path:", err)
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	logging.Log(true, "Absolute path:", absPath)

	// Launch mstsc.exe with the RDP file
	logging.Log(true, "Launching mstsc.exe with RDP file")
	cmd := exec.Command("mstsc", absPath)

	logging.Log(true, "Starting mstsc process...")
	if err := cmd.Start(); err != nil {
		logging.Log(true, "ERROR: Failed to launch mstsc:", err)
		return fmt.Errorf("failed to launch mstsc: %w", err)
	}

	logging.Log(true, "mstsc process started successfully, PID:", cmd.Process.Pid)
	return nil
}

// LaunchHost launches an RDP session for the specified host and user
func (g *Generator) LaunchHost(host models.Host, user models.User) error {
	logging.Log(true, "LaunchHost started for host:", host.Name, "address:", host.Address, "user:", user.Username)
	logging.Log(true, "User details - ID:", user.ID, "Name:", user.Name, "Username:", user.Username)
	logging.Log(true, "User has encrypted password:", user.EncryptedPassword != "")

	// User password should already be stored in Windows CredStore when user was saved
	if user.EncryptedPassword != "" {
		logging.Log(true, "User has encrypted password stored")
	}

	logging.Log(true, "Using existing credentials from Windows Credential Store (stored during user save)") // Generate standard RDP file (password comes from CredStore)
	logging.Log(true, "Generating RDP file")
	rdpFile, err := g.GenerateRDPFile(host, user)
	if err != nil {
		logging.Log(true, "ERROR: Failed to generate RDP file:", err)
		return fmt.Errorf("failed to generate RDP file: %w", err)
	}
	logging.Log(true, "RDP file generated:", rdpFile)

	// Launch RDP session
	logging.Log(true, "Launching RDP session")
	if err := g.LaunchRDP(rdpFile); err != nil {
		logging.Log(true, "ERROR: Failed to launch RDP session:", err)
		return fmt.Errorf("failed to launch RDP session: %w", err)
	}

	logging.Log(true, "LaunchHost completed successfully")
	return nil
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
