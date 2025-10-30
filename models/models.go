package models

import (
	"fmt"
	"time"
)

// User represents a user credential
type User struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`           // user@domain or email format
	Login             string    `json:"login"`              // actual login name
	Domain            string    `json:"domain"`             // domain name (optional)
	EncryptedPassword string    `json:"encrypted_password"` // AES encrypted password
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
}

// Host represents a remote host configuration
type Host struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"` // hostname or IP
	Port    int    `json:"port"`    // default 3389
	UserID  string `json:"user_id"` // reference to User.ID

	// RDP Settings
	RedirectClipboard bool   `json:"redirect_clipboard"`
	RedirectDrives    bool   `json:"redirect_drives"`
	DrivesToRedirect  string `json:"drives_to_redirect"` // "*" for all drives
	DisplayMode       string `json:"display_mode"`       // "fullscreen" or "window"
	DynamicResolution bool   `json:"dynamic_resolution"`
	ScreenMode        int    `json:"screen_mode"` // 1 = windowed, 2 = fullscreen

	// User-entered window size (what user actually wants)
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`

	// Calculated RDP desktop size (client area)
	DesktopWidth  int `json:"desktop_width"`
	DesktopHeight int `json:"desktop_height"`

	// Position and positioning
	PositionX int    `json:"position_x"`
	PositionY int    `json:"position_y"`
	WinPosStr string `json:"win_pos_str"` // calculated window position string

	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

// Users represents a collection of users
type Users struct {
	Users []User `json:"users"`
}

// Hosts represents a collection of hosts
type Hosts struct {
	Hosts []Host `json:"hosts"`
}

// NewUser creates a new user with generated ID and timestamps
func NewUser(name, username string) User {
	now := time.Now()
	return User{
		ID:         generateID(),
		Name:       name,
		Username:   username,
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

// NewHost creates a new host with generated ID and timestamps
func NewHost(name, address string, port int, userID string) Host {
	now := time.Now()
	return Host{
		ID:                generateID(),
		Name:              name,
		Address:           address,
		Port:              port,
		UserID:            userID,
		RedirectClipboard: true,
		RedirectDrives:    false,
		DrivesToRedirect:  "*",
		DisplayMode:       "window",
		DynamicResolution: true,
		ScreenMode:        1,    // windowed by default
		WindowWidth:       1200, // User-desired window size
		WindowHeight:      800,
		DesktopWidth:      1184, // Calculated client area (approx)
		DesktopHeight:     761,
		PositionX:         100,
		PositionY:         100,
		WinPosStr:         "0,1,100,100,1300,900", // Will be calculated based on position/size
		CreatedAt:         now,
		ModifiedAt:        now,
	}
}

// generateID generates a simple ID (you might want to use UUID in production)
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
