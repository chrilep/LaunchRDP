package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	AppName    = "Lancer"
	SubAppName = "LaunchRDP"
)

var (
	ConfigDir string
	DataDir   string
	LogsDir   string
	TempDir   string
)

// InitDirectories creates the necessary directories for the application
func InitDirectories() error {
	// Get base directories
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")

	if appData == "" || localAppData == "" {
		return fmt.Errorf("APPDATA or LOCALAPPDATA environment variables not found")
	}

	// Set directory paths
	ConfigDir = filepath.Join(appData, AppName, SubAppName)
	DataDir = ConfigDir
	LogsDir = filepath.Join(localAppData, AppName, SubAppName, "logs")
	TempDir = filepath.Join(localAppData, AppName, SubAppName, "temp")

	// Create directories
	dirs := []string{ConfigDir, DataDir, LogsDir, TempDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetConfigPath returns the full path to a config file
func GetConfigPath(filename string) string {
	return filepath.Join(ConfigDir, filename)
}

// GetTempPath returns the full path to a temporary file
func GetTempPath(filename string) string {
	return filepath.Join(TempDir, filename)
}

// GetLogPath returns the full path to a log file
func GetLogPath(filename string) string {
	return filepath.Join(LogsDir, filename)
}
