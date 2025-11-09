package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/chrilep/LaunchRDP/app/config"
	"github.com/chrilep/LaunchRDP/app/credentials"
	"github.com/chrilep/LaunchRDP/app/logging"
	"github.com/chrilep/LaunchRDP/app/models"
	"github.com/chrilep/LaunchRDP/app/rdp"
	"github.com/chrilep/LaunchRDP/app/storage"
)

// LaunchRDPApp manages storage, credentials, RDP generator and window state
type LaunchRDPApp struct {
	ctx              context.Context
	storage          *storage.Storage
	credManager      *credentials.CredentialManager
	rdpGen           *rdp.Generator
	winState         *WindowState
	winStateMu       sync.Mutex
	windowHWND       uintptr
	winEventHook     uintptr
	winEventCallback uintptr
	stopTicker       chan struct{}
}

// NewLaunchRDPApp erstellt die App mit Default-WindowState (intended -7,0)
func NewLaunchRDPApp() *LaunchRDPApp {
	debug := true
	logging.Log(debug, "Creating new LaunchRDP Wails app instance (simplified)")
	app := &LaunchRDPApp{
		storage:     storage.NewStorage(),
		credManager: credentials.NewCredentialManager(),
		rdpGen:      rdp.NewGenerator(),
		winState:    &WindowState{X: -7, Y: 0, Width: 275, Height: 500, DeltaX: 0, DeltaY: 0},
	}
	app.rdpGen.SetSaveUserCallback(app.saveUserAfterMigration)
	return app
}

// Startup loads saved window state
func (a *LaunchRDPApp) Startup(ctx context.Context) {
	a.ctx = ctx
	if ws, err := loadWindowState(); err == nil && ws != nil {
		a.winState = ws
		logging.Log(true, fmt.Sprintf("Loaded window state physical(%d,%d %dx%d) delta(%d,%d)", ws.X, ws.Y, ws.Width, ws.Height, ws.DeltaX, ws.DeltaY))
	} else {
		logging.Log(true, "Using default window state")
	}
}

// DomReady: set intended position (-7,0 minus stored delta) then record external shift (e.g. DockFinder 30px)
func (a *LaunchRDPApp) DomReady(ctx context.Context) {
	debug := true
	intendedX := a.winState.X - a.winState.DeltaX
	intendedY := a.winState.Y - a.winState.DeltaY
	logging.Log(debug, fmt.Sprintf("DomReady: intendedX=%d intendedY=%d physicalX=%d physicalY=%d deltaX=%d deltaY=%d", intendedX, intendedY, a.winState.X, a.winState.Y, a.winState.DeltaX, a.winState.DeltaY))
	runtime.WindowSetSize(ctx, a.winState.Width, a.winState.Height)
	runtime.WindowSetPosition(ctx, intendedX, intendedY)
	runtime.WindowShow(ctx)
	fx, fy := runtime.WindowGetPosition(ctx)
	if fx != intendedX || fy != intendedY {
		a.winState.DeltaX = fx - intendedX
		a.winState.DeltaY = fy - intendedY
		logging.Log(debug, fmt.Sprintf("DomReady: external shift detected finalX=%d finalY=%d deltaX=%d deltaY=%d", fx, fy, a.winState.DeltaX, a.winState.DeltaY))
	} else {
		logging.Log(debug, fmt.Sprintf("DomReady: no external shift physicalX=%d physicalY=%d deltaX=%d deltaY=%d", fx, fy, a.winState.DeltaX, a.winState.DeltaY))
	}
	a.winState.X = fx
	a.winState.Y = fy
	_ = saveWindowState(a.winState) // persist early so next start has physical+delta

	// Start proactive geometry monitoring (system messages + ticker)
	a.initHWND()
	if err := a.startMoveSizeHook(); err != nil {
		logging.Log(true, "DomReady: startMoveSizeHook failed", err)
	}
	a.startGeometryTicker()
}

// CreateUser creates a new user and saves it
func (a *LaunchRDPApp) CreateUser(username, login, domain, password string) error {
	debug := false
	logging.Log(debug, "API: Creating user -", username, "Domain:", domain)
	user := models.NewUser(username, username)
	user.Login = login
	user.Domain = domain

	// Store encrypted password as backup (DPAPI)
	// We'll only decrypt it when assigning to hosts or updating credentials
	if password != "" {
		encryptedPassword, err := a.credManager.EncryptPasswordDPAPI(password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = encryptedPassword
		logging.Log(debug, "Password encrypted with DPAPI for user:", username)
	}

	users, _ := a.storage.LoadUsers()
	users = append(users, user)
	if err := a.storage.SaveUsers(users); err != nil {
		return err
	}
	return nil
}

// GetUsers returns all saved users
func (a *LaunchRDPApp) GetUsers() ([]models.User, error) {
	debug := false
	logging.Log(debug, "API: Loading users")
	users, err := a.storage.LoadUsers()
	if err != nil {
		logging.Log(true, "ERROR: Failed to load users:", err)
		return nil, err
	}
	logging.Log(debug, "API: Loaded", len(users), "users")
	return users, nil
}

// UpdateUser updates basic data + password (if not __UNCHANGED__)
func (a *LaunchRDPApp) UpdateUser(userID, username, login, domain, password string) error {
	debug := false
	users, err := a.storage.LoadUsers()
	if err != nil {
		return err
	}
	idx := -1
	for i, u := range users {
		if u.ID == userID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("user not found")
	}
	usr := &users[idx]
	usr.Username = username
	usr.Login = login
	usr.Domain = domain

	if password != "" && password != "__UNCHANGED__" {
		enc, err := a.credManager.EncryptPasswordDPAPI(password)
		if err != nil {
			return err
		}
		usr.EncryptedPassword = enc
		hosts, _ := a.storage.LoadHosts()
		logging.Log(true, "UpdateUser: Storing credentials for hosts associated with user", username)
		for _, h := range hosts {
			if h.UserID == userID {
				logging.Log(true, "UpdateUser: Calling StoreCredential for host:", h.Address, "user:", username)
				err := a.credManager.StoreCredential(h.Address, username, password)
				if err != nil {
					logging.Log(true, "UpdateUser: ERROR storing credential for host", h.Address, ":", err)
				} else {
					logging.Log(true, "UpdateUser: Successfully stored credential for host", h.Address)
				}
			}
		}
	}
	if err := a.storage.SaveUsers(users); err != nil {
		return err
	}
	logging.Log(debug, "API: User updated", userID)
	return nil
}

// DeleteUser - Replaces DELETE /api/users/{id}
func (a *LaunchRDPApp) DeleteUser(userID string) error {
	debug := false
	logging.Log(debug, "API: Deleting user", userID)

	// Load users and filter out the deleted one
	users, err := a.storage.LoadUsers()
	if err != nil {
		return err
	}

	var updatedUsers []models.User
	var deletedUser *models.User

	for _, user := range users {
		if user.ID == userID {
			deletedUser = &user
		} else {
			updatedUsers = append(updatedUsers, user)
		}
	}

	if deletedUser == nil {
		logging.Log(true, "ERROR: User not found:", userID)
		return fmt.Errorf("user not found")
	}

	// Delete credentials
	hosts, _ := a.storage.LoadHosts()
	for _, host := range hosts {
		if host.UserID == userID {
			a.credManager.DeleteCredential(host.Address)
		}
	}

	err = a.storage.SaveUsers(updatedUsers)
	if err != nil {
		logging.Log(true, "ERROR: Failed to delete user:", err)
		return err
	}

	logging.Log(debug, "User deleted successfully:", userID)
	return nil
}

// GetHosts - Replaces GET /api/hosts
func (a *LaunchRDPApp) GetHosts() ([]models.Host, error) {
	debug := false
	logging.Log(debug, "API: Loading hosts")
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		logging.Log(true, "ERROR: Failed to load hosts:", err)
		return nil, err
	}
	logging.Log(debug, "API: Loaded", len(hosts), "hosts")
	return hosts, nil
}

// CreateHost - Replaces POST /api/hosts
// Extended: Additional optional parameters for RDP settings
func (a *LaunchRDPApp) CreateHost(name, address, userID string, port int) error {
	debug := false
	logging.Log(debug, "API: Creating host -", name, "Address:", address, "Port:", port)

	// Create new host - same logic as server.go
	host := models.NewHost(name, address, port, userID)
	// Note: Extended storage is handled via separate Update function after creation if needed.

	// Save host using array pattern
	hosts, _ := a.storage.LoadHosts()
	hosts = append(hosts, host)
	err := a.storage.SaveHosts(hosts)
	if err != nil {
		logging.Log(true, "ERROR: Failed to save host:", err)
		return err
	}

	logging.Log(debug, "Host created successfully:", name, "with ID:", host.ID)
	return nil
}

// UpdateHost - basic update (name/address/port/user). Advanced fields use UpdateHostFull.
func (a *LaunchRDPApp) UpdateHost(hostID, name, address, userID string, port int) error {
	debug := false
	logging.Log(debug, "API: Updating host", hostID)

	// Load hosts and find the one to update
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		return err
	}

	hostIndex := -1
	var oldAddress string
	var oldUserID string
	for i, host := range hosts {
		if host.ID == hostID {
			hostIndex = i
			oldAddress = host.Address
			oldUserID = host.UserID
			break
		}
	}

	if hostIndex == -1 {
		logging.Log(true, "ERROR: Host not found:", hostID)
		return fmt.Errorf("host not found")
	}

	// Basic host data update
	h := &hosts[hostIndex]
	h.Name = name
	h.Address = address
	h.Port = port
	h.UserID = userID

	err = a.storage.SaveHosts(hosts)
	if err != nil {
		logging.Log(true, "ERROR: Failed to update host:", err)
		return err
	}

	// Update credentials if address or user changed
	if address != oldAddress || userID != oldUserID {
		logging.Log(debug, "Host address or user changed, updating credentials")

		// Delete old credential if address changed
		if address != oldAddress {
			logging.Log(debug, "Deleting old credential for:", oldAddress)
			a.credManager.DeleteCredential(oldAddress)
		}

		// Store new credential
		users, err := a.storage.LoadUsers()
		if err == nil {
			for _, user := range users {
				if user.ID == userID {
					if user.EncryptedPassword != "" {
						password, err := a.credManager.DecryptPasswordDPAPI(user.EncryptedPassword)
						if err == nil {
							logging.Log(debug, "Storing credential for new host/user:", address, user.Username)
							err = a.credManager.StoreCredential(address, user.Username, password)
							if err != nil {
								logging.Log(true, "ERROR: Failed to store credential:", err)
							}
						}
					}
					break
				}
			}
		}
	}

	logging.Log(debug, "Host updated successfully:", name)
	return nil
}

// UpdateHostFull updates all host settings including RDP/window properties.
func (a *LaunchRDPApp) UpdateHostFull(hostID, name, address, userID string, port int,
	displayMode string, positionX, positionY, windowWidth, windowHeight int,
	redirectClipboard, redirectDrives bool) error {
	debug := false
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		return err
	}
	idx := -1
	var oldAddress string
	var oldUserID string
	for i, host := range hosts {
		if host.ID == hostID {
			idx = i
			oldAddress = host.Address
			oldUserID = host.UserID
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("host not found")
	}
	h := &hosts[idx]
	h.Name = name
	h.Address = address
	h.Port = port
	h.UserID = userID
	h.DisplayMode = displayMode
	if displayMode == "fullscreen" {
		h.ScreenMode = 2
	} else {
		h.ScreenMode = 1
	}
	h.PositionX = positionX
	h.PositionY = positionY
	h.WindowWidth = windowWidth
	h.WindowHeight = windowHeight
	// Approximate client area
	h.DesktopWidth = windowWidth - 16
	h.DesktopHeight = windowHeight - 59
	h.RedirectClipboard = redirectClipboard
	h.RedirectDrives = redirectDrives
	h.ModifiedAt = time.Now()
	if err := a.storage.SaveHosts(hosts); err != nil {
		return err
	}

	// Update credentials if address or user changed
	if address != oldAddress || userID != oldUserID {
		logging.Log(debug, "Host address or user changed, updating credentials")

		// Delete old credential if address changed
		if address != oldAddress {
			logging.Log(debug, "Deleting old credential for:", oldAddress)
			a.credManager.DeleteCredential(oldAddress)
		}

		// Store new credential
		users, err := a.storage.LoadUsers()
		if err == nil {
			for _, user := range users {
				if user.ID == userID {
					if user.EncryptedPassword != "" {
						password, err := a.credManager.DecryptPasswordDPAPI(user.EncryptedPassword)
						if err == nil {
							logging.Log(debug, "Storing credential for new host/user:", address, user.Username)
							err = a.credManager.StoreCredential(address, user.Username, password)
							if err != nil {
								logging.Log(true, "ERROR: Failed to store credential:", err)
							}
						}
					}
					break
				}
			}
		}
	}

	return nil
}

// CreateHostFull creates a host with advanced properties.
func (a *LaunchRDPApp) CreateHostFull(name, address, userID string, port int,
	displayMode string, positionX, positionY, windowWidth, windowHeight int,
	redirectClipboard, redirectDrives bool) error {
	host := models.NewHost(name, address, port, userID)
	host.DisplayMode = displayMode
	if displayMode == "fullscreen" {
		host.ScreenMode = 2
	} else {
		host.ScreenMode = 1
	}
	host.PositionX = positionX
	host.PositionY = positionY
	host.WindowWidth = windowWidth
	host.WindowHeight = windowHeight
	host.DesktopWidth = windowWidth - 16
	host.DesktopHeight = windowHeight - 59
	host.RedirectClipboard = redirectClipboard
	host.RedirectDrives = redirectDrives
	hosts, _ := a.storage.LoadHosts()
	hosts = append(hosts, host)
	if err := a.storage.SaveHosts(hosts); err != nil {
		return err
	}
	return nil
}

// GenerateHostRDP creates/updates the RDP file for a given host (used after save)
func (a *LaunchRDPApp) GenerateHostRDP(hostID string) (string, error) {
	debug := false
	logging.Log(debug, "API: GenerateHostRDP invoked for host", hostID)
	// Load host
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		return "", err
	}
	var host *models.Host
	for i := range hosts {
		if hosts[i].ID == hostID {
			host = &hosts[i]
			break
		}
	}
	if host == nil {
		return "", fmt.Errorf("host not found")
	}
	if host.UserID == "" {
		return "", fmt.Errorf("host has no assigned user")
	}

	// Load user from storage
	users, err := a.storage.LoadUsers()
	if err != nil {
		return "", err
	}

	var user *models.User
	for i := range users {
		if users[i].ID == host.UserID {
			user = &users[i]
			break
		}
	}
	if user == nil {
		return "", fmt.Errorf("assigned user not found")
	}

	// Generate RDP file
	path, err := a.rdpGen.GenerateRDPFile(*host, *user)
	if err != nil {
		return "", err
	}
	logging.Log(debug, "RDP file generated at", path)
	return path, nil
}

// DeleteHost - Replaces DELETE /api/hosts/{id}
func (a *LaunchRDPApp) DeleteHost(hostID string) error {
	debug := false
	logging.Log(debug, "API: Deleting host", hostID)

	// Load hosts and filter out the deleted one
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		return err
	}

	var updatedHosts []models.Host
	found := false

	for _, host := range hosts {
		if host.ID == hostID {
			found = true
		} else {
			updatedHosts = append(updatedHosts, host)
		}
	}

	if !found {
		logging.Log(true, "ERROR: Host not found:", hostID)
		return fmt.Errorf("host not found")
	}

	err = a.storage.SaveHosts(updatedHosts)
	if err != nil {
		logging.Log(true, "ERROR: Failed to delete host:", err)
		return err
	}

	logging.Log(debug, "Host deleted successfully:", hostID)
	return nil
}

// MousePosition represents a screen position
type MousePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GetMousePosition - Replaces GET /api/mouse-position
func (a *LaunchRDPApp) GetMousePosition() (*MousePosition, error) {
	debug := false
	logging.Log(debug, "API: Getting mouse position")

	user32 := syscall.NewLazyDLL("user32.dll")
	procGetCursorPos := user32.NewProc("GetCursorPos")
	type point struct {
		X int32
		Y int32
	}
	var pt point
	r, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if r == 0 { // failed
		logging.Log(true, "ERROR: GetCursorPos failed", err)
		return &MousePosition{X: 0, Y: 0}, nil
	}
	position := &MousePosition{X: int(pt.X), Y: int(pt.Y)}
	logging.Log(debug, "Mouse position:", position.X, position.Y)
	return position, nil
}

// LaunchRDP - Replaces POST /api/launch - THE MAIN FUNCTION!
func (a *LaunchRDPApp) LaunchRDP(hostID, userID string, positionX, positionY int32) (bool, error) {
	debug := false
	logging.Log(debug, "API: Launching RDP - Host:", hostID, "User:", userID, "Position:", positionX, positionY)

	// Create position struct
	_ = &MousePosition{X: int(positionX), Y: int(positionY)} // Position not used by LaunchHost

	// Load host
	hosts, err := a.storage.LoadHosts()
	if err != nil {
		logging.Log(true, "ERROR: Failed to load hosts:", err)
		return false, err
	}

	var host *models.Host
	for _, h := range hosts {
		if h.ID == hostID {
			host = &h
			break
		}
	}

	if host == nil {
		logging.Log(true, "ERROR: Host not found:", hostID)
		return false, fmt.Errorf("host not found")
	}
	logging.Log(debug, "Host loaded:", host.Name)

	// Load user - handle system user specially
	var user *models.User

	// Load user from storage
	users, err := a.storage.LoadUsers()
	if err != nil {
		logging.Log(true, "ERROR: Failed to load users:", err)
		return false, err
	}

	for _, u := range users {
		if u.ID == userID {
			user = &u
			break
		}
	}

	if user == nil {
		logging.Log(true, "ERROR: User not found:", userID)
		return false, fmt.Errorf("user not found")
	}
	logging.Log(debug, "User loaded:", user.Username)

	// Password is already stored in Windows Credential Manager
	// No need to decrypt from DPAPI here - Windows will handle it automatically
	logging.Log(debug, "Credentials will be loaded from Windows Credential Manager")

	// Generate and launch RDP connection
	wasReused, err := a.rdpGen.LaunchHost(*host, *user)
	if err != nil {
		logging.Log(true, "ERROR: Failed to launch RDP:", err)
		return false, err
	}

	if wasReused {
		logging.Log(debug, "RDP window reused (existing connection activated)")
	} else {
		logging.Log(debug, "RDP connection launched successfully!")
	}
	return wasReused, nil
}

// WindowBorderInfo represents window border information
type WindowBorderInfo struct {
	Left         int `json:"left"`
	Right        int `json:"right"`
	Top          int `json:"top"`
	Bottom       int `json:"bottom"`
	ClientWidth  int `json:"clientWidth"`
	ClientHeight int `json:"clientHeight"`
	WindowWidth  int `json:"windowWidth"`
	WindowHeight int `json:"windowHeight"`
}

// GetWindowInfo - Replaces GET /api/window-info
func (a *LaunchRDPApp) GetWindowInfo() (*WindowBorderInfo, error) {
	debug := false
	logging.Log(debug, "API: Getting window info - not implemented in Wails yet")

	// For now, return a placeholder - can be implemented later if needed
	info := &WindowBorderInfo{
		Left:         0,
		Right:        0,
		Top:          30, // MyDockFinder offset
		Bottom:       0,
		ClientWidth:  275,
		ClientHeight: 1349,
		WindowWidth:  275,
		WindowHeight: 1349,
	}

	logging.Log(debug, "Window info placeholder returned")
	return info, nil
} // LogMessage - For frontend logging - Replaces POST /api/log
func (a *LaunchRDPApp) LogMessage(level, message string) {
	isError := level == "error"
	logging.Log(isError, "Frontend:", message)
}

// GetWindowBorderInfo returns the real window frame metrics from the OS so the
// frontend can calculate the usable client size and WinPosStr accurately.
// This uses GetSystemMetrics for size frame, padded border and caption height.
// Reference metrics:
//
//	SM_CXSIZEFRAME / SM_CYSIZEFRAME: Thickness of the sizing border
//	SM_CXPADDEDBORDER: Extra border padding introduced with Aero
//	SM_CYCAPTION: Caption (title bar) height
//
// We return a simplified model: left/right/top/bottom total non-client offsets.
func (a *LaunchRDPApp) GetWindowBorderInfo() (*WindowBorderInfo, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	procGetSystemMetrics := user32.NewProc("GetSystemMetrics")

	// Helper to read metric
	getMetric := func(index int) int {
		r, _, _ := procGetSystemMetrics.Call(uintptr(index))
		return int(r)
	}

	const (
		SM_CXSIZEFRAME    = 32
		SM_CYSIZEFRAME    = 33
		SM_CXPADDEDBORDER = 92
		SM_CYCAPTION      = 4
	)

	sizeFrameX := getMetric(SM_CXSIZEFRAME)
	sizeFrameY := getMetric(SM_CYSIZEFRAME)
	padded := getMetric(SM_CXPADDEDBORDER)
	caption := getMetric(SM_CYCAPTION)

	// Non-client borders left/right include size frame + padded border
	left := sizeFrameX + padded
	right := sizeFrameX + padded
	top := sizeFrameY + padded + caption
	bottom := sizeFrameY + padded

	// We do not know the runtime window size here; fill with zeroes and let frontend compute client size from user input.
	info := &WindowBorderInfo{
		Left:         left,
		Right:        right,
		Top:          top,
		Bottom:       bottom,
		ClientWidth:  0,
		ClientHeight: 0,
		WindowWidth:  0,
		WindowHeight: 0,
	}
	return info, nil
}

// ================= Window State Persistence =================
// WindowState holds last window geometry
type WindowState struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	DeltaX int `json:"deltaX"`
	DeltaY int `json:"deltaY"`
}

// captureCurrentWindow safely reads current window position & size.
// It wraps runtime calls with recover because during shutdown the underlying
// native window/control may already be tearing down (causing internal panic like divide by zero in DPI scale).
// Returns (x,y,w,h) and a bool indicating success.
// (former captureCurrentWindow removed; proactive hooks now maintain state)

// initHWND finds the window handle by title (class optional) after window is shown
func (a *LaunchRDPApp) initHWND() {
	if a.windowHWND != 0 {
		return
	}
	user32 := syscall.NewLazyDLL("user32.dll")
	procFindWindow := user32.NewProc("FindWindowW")
	titlePtr, _ := syscall.UTF16PtrFromString("LaunchRDP")
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	a.windowHWND = hwnd
	logging.Log(true, fmt.Sprintf("initHWND: hwnd=%x", hwnd))
}

// startGeometryTicker periodically samples window geometry as a fallback
func (a *LaunchRDPApp) startGeometryTicker() {
	if a.stopTicker != nil {
		return
	}
	a.stopTicker = make(chan struct{})
	ticker := time.NewTicker(5000 * time.Millisecond) // 5 seconds - hook handles most updates
	go func() {
		debug := false // Only log when changes detected
		for {
			select {
			case <-a.stopTicker:
				ticker.Stop()
				return
			case <-ticker.C:
				if a.ctx == nil {
					continue
				}
				x, y := runtime.WindowGetPosition(a.ctx)
				w, h := runtime.WindowGetSize(a.ctx)
				a.winStateMu.Lock()
				changed := (x != a.winState.X || y != a.winState.Y || w != a.winState.Width || h != a.winState.Height)
				if changed {
					a.winState.X = x
					a.winState.Y = y
					a.winState.Width = w
					a.winState.Height = h
					if err := saveWindowState(a.winState); err == nil {
						logging.Log(debug, fmt.Sprintf("Ticker: updated & saved X=%d Y=%d W=%d H=%d", x, y, w, h))
					}
				}
				a.winStateMu.Unlock()
			}
		}
	}()
}

// startMoveSizeHook sets a WinEvent hook for MOVESIZEEND events to capture final geometry after user interaction
func (a *LaunchRDPApp) startMoveSizeHook() error {
	user32 := syscall.NewLazyDLL("user32.dll")
	procSetWinEventHook := user32.NewProc("SetWinEventHook")
	procUnhookWinEvent := user32.NewProc("UnhookWinEvent")
	procGetWindowRect := user32.NewProc("GetWindowRect")
	procGetWindowThreadProcessId := user32.NewProc("GetWindowThreadProcessId")

	const EVENT_SYSTEM_MOVESIZEEND = 0x000B
	const WINEVENT_INCONTEXT = 0x0004

	if a.windowHWND == 0 {
		a.initHWND()
	}
	if a.windowHWND == 0 {
		return fmt.Errorf("startMoveSizeHook: hwnd not found")
	}

	// Get the thread ID of the window to hook events for that thread only
	var processID uint32
	threadID, _, _ := procGetWindowThreadProcessId.Call(a.windowHWND, uintptr(unsafe.Pointer(&processID)))
	logging.Log(true, fmt.Sprintf("startMoveSizeHook: hwnd=%x threadId=%d processId=%d", a.windowHWND, threadID, processID))

	// Load the current module handle for INCONTEXT hook
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")
	hModule, _, _ := procGetModuleHandle.Call(0)

	// Callback - INCONTEXT runs in the UI thread's message loop (compatible with Wails)
	a.winEventCallback = syscall.NewCallback(func(hWinEventHook, event, hwnd, idObject, idChild, dwEventThread, dwmsEventTime uintptr) uintptr {
		if event == EVENT_SYSTEM_MOVESIZEEND && hwnd == a.windowHWND {
			type rect struct{ Left, Top, Right, Bottom int32 }
			var r rect
			rOk, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
			if rOk != 0 {
				x := int(r.Left)
				y := int(r.Top)
				w := int(r.Right - r.Left)
				h := int(r.Bottom - r.Top)
				a.winStateMu.Lock()
				a.winState.X = x
				a.winState.Y = y
				a.winState.Width = w
				a.winState.Height = h
				_ = saveWindowState(a.winState)
				a.winStateMu.Unlock()
				logging.Log(true, fmt.Sprintf("Hook MOVESIZEEND: saved X=%d Y=%d W=%d H=%d", x, y, w, h))
			}
		}
		return 0
	})

	hHook, _, err := procSetWinEventHook.Call(
		uintptr(EVENT_SYSTEM_MOVESIZEEND), // eventMin
		uintptr(EVENT_SYSTEM_MOVESIZEEND), // eventMax
		hModule,                           // hmod - our module for INCONTEXT
		a.winEventCallback,                // callback
		uintptr(processID),                // idProcess - our process only
		threadID,                          // idThread - our UI thread only
		uintptr(WINEVENT_INCONTEXT),       // flags - INCONTEXT for message loop integration
	)
	if hHook == 0 {
		return fmt.Errorf("SetWinEventHook failed: %v", err)
	}
	a.winEventHook = hHook
	_ = procUnhookWinEvent // silence unused (we unhook in Shutdown)
	logging.Log(true, fmt.Sprintf("startMoveSizeHook: hook=%x", hHook))
	return nil
}

// WorkArea represents usable desktop area excluding taskbar
type WorkArea struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// MonitorWorkArea enthält vollständige Monitor- und WorkArea-Daten für jeden Bildschirm
type MonitorWorkArea struct {
	Index         int  `json:"index"`
	MonitorLeft   int  `json:"monitorLeft"`
	MonitorTop    int  `json:"monitorTop"`
	MonitorRight  int  `json:"monitorRight"`
	MonitorBottom int  `json:"monitorBottom"`
	WorkLeft      int  `json:"workLeft"`
	WorkTop       int  `json:"workTop"`
	WorkRight     int  `json:"workRight"`
	WorkBottom    int  `json:"workBottom"`
	Primary       bool `json:"primary"`
}

const WindowFileName = "window.json"

func loadWindowState() (*WindowState, error) {
	path := config.GetConfigPath(WindowFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			def := &WindowState{X: -7, Y: 0, Width: 275, Height: 500, DeltaX: 0, DeltaY: 0}
			if b, mErr := json.MarshalIndent(def, "", "  "); mErr == nil {
				_ = os.WriteFile(path, b, 0644)
			}
			return def, nil
		}
		return nil, err
	}
	var ws WindowState
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	if ws.Width < 50 {
		ws.Width = 275
	}
	if ws.Height < 50 {
		ws.Height = 500
	}
	return &ws, nil
}

func saveWindowState(ws *WindowState) error {
	if ws == nil {
		return fmt.Errorf("nil window state")
	}
	path := config.GetConfigPath(WindowFileName)
	b, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

// GetWorkArea retrieves primary monitor work area (excluding taskbar)
func (a *LaunchRDPApp) GetWorkArea() (*WorkArea, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	procSPI := user32.NewProc("SystemParametersInfoW")
	const SPI_GETWORKAREA = 0x0030
	type rect struct{ Left, Top, Right, Bottom int32 }
	var r rect
	ret, _, err := procSPI.Call(uintptr(SPI_GETWORKAREA), 0, uintptr(unsafe.Pointer(&r)), 0)
	if ret == 0 {
		return nil, err
	}
	wa := &WorkArea{Left: int(r.Left), Top: int(r.Top), Right: int(r.Right), Bottom: int(r.Bottom)}
	wa.Width = wa.Right - wa.Left
	wa.Height = wa.Bottom - wa.Top
	return wa, nil
}

// GetMonitorWorkAreas enumeriert alle Monitore (EnumDisplayMonitors) und liest deren rcMonitor & rcWork
func (a *LaunchRDPApp) GetMonitorWorkAreas() ([]MonitorWorkArea, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	procEnum := user32.NewProc("EnumDisplayMonitors")
	procGetInfo := user32.NewProc("GetMonitorInfoW")
	debug := true
	logging.Log(debug, "API: GetMonitorWorkAreas start")

	type rect struct{ Left, Top, Right, Bottom int32 }
	type monitorInfoEx struct {
		cbSize    uint32
		rcMonitor rect
		rcWork    rect
		dwFlags   uint32
		szDevice  [32]uint16
	}

	monitors := make([]MonitorWorkArea, 0, 8)
	index := 0

	enumCallback := syscall.NewCallback(func(hMonitor, hdc, lprcMonitor, dwData uintptr) uintptr {
		var mi monitorInfoEx
		mi.cbSize = uint32(unsafe.Sizeof(mi))
		r, _, _ := procGetInfo.Call(hMonitor, uintptr(unsafe.Pointer(&mi)))
		if r != 0 {
			m := MonitorWorkArea{
				Index:         index,
				MonitorLeft:   int(mi.rcMonitor.Left),
				MonitorTop:    int(mi.rcMonitor.Top),
				MonitorRight:  int(mi.rcMonitor.Right),
				MonitorBottom: int(mi.rcMonitor.Bottom),
				WorkLeft:      int(mi.rcWork.Left),
				WorkTop:       int(mi.rcWork.Top),
				WorkRight:     int(mi.rcWork.Right),
				WorkBottom:    int(mi.rcWork.Bottom),
				Primary:       (mi.dwFlags & 1) == 1, // MONITORINFOF_PRIMARY
			}
			monitors = append(monitors, m)
			index++
		}
		return 1 // continue enumeration
	})

	r, _, err := procEnum.Call(0, 0, enumCallback, 0)
	if r == 0 {
		logging.Log(true, "API: GetMonitorWorkAreas failed", err)
		return nil, err
	}
	logging.Log(debug, fmt.Sprintf("API: GetMonitorWorkAreas success monitors=%d", len(monitors)))
	return monitors, nil
}

// GetWindowState returns current window geometry + detected offsets
func (a *LaunchRDPApp) GetWindowState() (*WindowState, error) {
	debug := true
	logging.Log(debug, "API: GetWindowState start (simplified)")
	if a.ctx == nil {
		return a.winState, nil
	}
	x, y := runtime.WindowGetPosition(a.ctx)
	w, h := runtime.WindowGetSize(a.ctx)
	a.winState.X = x
	a.winState.Y = y
	a.winState.Width = w
	a.winState.Height = h
	logging.Log(debug, fmt.Sprintf("GetWindowState: physicalX=%d physicalY=%d width=%d height=%d deltaX=%d deltaY=%d", x, y, w, h, a.winState.DeltaX, a.winState.DeltaY))
	return a.winState, nil
}

// SaveWindowGeometry liest aktuelle Position/Größe und persistiert sie.
// Wird von Shutdown aufgerufen; kann optional auch vom Frontend genutzt werden.
// PersistWindowState writes current in-memory geometry without runtime access (safe during teardown)
func (a *LaunchRDPApp) PersistWindowState() error {
	a.winStateMu.Lock()
	defer a.winStateMu.Unlock()
	if err := saveWindowState(a.winState); err != nil {
		logging.Log(true, "PersistWindowState: save failed", err)
		return err
	}
	logging.Log(true, fmt.Sprintf("PersistWindowState: saved X=%d Y=%d W=%d H=%d", a.winState.X, a.winState.Y, a.winState.Width, a.winState.Height))
	return nil
}

// Shutdown is called when the app terminates; persist window geometry
func (a *LaunchRDPApp) Shutdown(ctx context.Context) {
	// Stop ticker
	if a.stopTicker != nil {
		close(a.stopTicker)
	}
	// Unhook win event if set
	if a.winEventHook != 0 {
		user32 := syscall.NewLazyDLL("user32.dll")
		procUnhookWinEvent := user32.NewProc("UnhookWinEvent")
		_, _, _ = procUnhookWinEvent.Call(a.winEventHook)
		logging.Log(true, fmt.Sprintf("Shutdown: unhooked winEventHook=%x", a.winEventHook))
	}
	// Persist in-memory state only (no runtime calls – window may be gone)
	_ = a.PersistWindowState()
}

// saveUserAfterMigration - Callback for RDP generator - SAME AS BEFORE!
func (a *LaunchRDPApp) saveUserAfterMigration(user models.User) error {
	debug := false
	logging.Log(debug, "Saving user after migration:", user.Username)

	// Load users, add/update user, save back
	users, _ := a.storage.LoadUsers()

	// Check if user exists
	userIndex := -1
	for i, u := range users {
		if u.ID == user.ID {
			userIndex = i
			break
		}
	}

	if userIndex >= 0 {
		// Update existing user
		users[userIndex] = user
	} else {
		// Add new user
		users = append(users, user)
	}

	return a.storage.SaveUsers(users)
}
