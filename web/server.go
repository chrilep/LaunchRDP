package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/chrilep/LaunchRDP/credentials"
	"github.com/chrilep/LaunchRDP/logging"
	"github.com/chrilep/LaunchRDP/models"
	"github.com/chrilep/LaunchRDP/rdp"
	"github.com/chrilep/LaunchRDP/storage"
)

// Server represents the web server
type Server struct {
	storage      *storage.Storage
	credManager  *credentials.CredentialManager
	rdpGenerator *rdp.Generator
	port         int
}

// NewServer creates a new web server
func NewServer(port int) (*Server, error) {
	server := &Server{
		storage:      storage.NewStorage(),
		credManager:  credentials.NewCredentialManager(),
		rdpGenerator: rdp.NewGenerator(),
		port:         port,
	}

	// Set up RDP generator callback
	server.rdpGenerator.SetSaveUserCallback(server.saveUserAfterMigration)

	return server, nil
}

// Start starts the web server
func (s *Server) Start() error {
	debug := true
	logging.Log(debug, "Starting web server on port", s.port)

	// Create new UI instance
	ui, err := NewUI()
	if err != nil {
		return err
	}

	// Register routes
	mux := http.NewServeMux()
	mux.Handle("/", ui)
	mux.HandleFunc("/api/users", s.handleUsers)
	mux.HandleFunc("/api/users/", s.handleUserByID)
	mux.HandleFunc("/api/hosts", s.handleHosts)
	mux.HandleFunc("/api/hosts/", s.handleHostByID)
	mux.HandleFunc("/api/launch", s.handleLaunch)
	mux.HandleFunc("/api/window-info", s.handleWindowInfo)
	mux.HandleFunc("/api/log", s.handleWebViewLog)

	// Start server - bind only to localhost to avoid firewall prompts
	addr := fmt.Sprintf("localhost:%d", s.port)
	logging.Log(debug, "Web server listening on", addr)

	// No browser opening - WebView2 only!
	return http.ListenAndServe(addr, mux)
}

// API handlers
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	debug := true
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		users, err := s.storage.LoadUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)

	case "POST":
		logging.Log(debug, "API: Creating new user")
		var userData struct {
			Username string `json:"username"`
			Login    string `json:"login"`
			Domain   string `json:"domain"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			logging.Log(true, "ERROR: Failed to decode user data:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logging.Log(debug, "User data received - Username:", userData.Username, "Domain:", userData.Domain)

		// Create new user
		user := models.NewUser(userData.Username, userData.Username)
		user.Login = userData.Login
		user.Domain = userData.Domain

		// Encrypt password
		if userData.Password != "" {
			logging.Log(debug, "Encrypting password with DPAPI")
			encryptedPassword, err := s.credManager.EncryptPasswordDPAPI(userData.Password)
			if err != nil {
				logging.Log(true, "ERROR: Failed to encrypt password:", err)
				http.Error(w, "Failed to encrypt password: "+err.Error(), http.StatusInternalServerError)
				return
			}
			user.EncryptedPassword = encryptedPassword
			logging.Log(debug, "Password encrypted successfully")

			// Store in Windows Credential Store
			logging.Log(debug, "Storing credentials in Windows Credential Store")
			hosts, _ := s.storage.LoadHosts()
			for _, host := range hosts {
				if host.UserID == user.ID {
					logging.Log(debug, "Storing credential for host:", host.Address, "user:", user.Username)
					s.credManager.StoreCredential(host.Address, user.Username, userData.Password)
				}
			}
		}

		// Save user
		logging.Log(debug, "Saving user to storage")
		users, _ := s.storage.LoadUsers()
		users = append(users, user)
		if err := s.storage.SaveUsers(users); err != nil {
			logging.Log(true, "ERROR: Failed to save users:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logging.Log(debug, "User saved successfully with ID:", user.ID)

		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	debug := false

	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from URL
	userID := filepath.Base(r.URL.Path)

	switch r.Method {
	case "PUT":
		var userData struct {
			Username string `json:"username"`
			Login    string `json:"login"`
			Domain   string `json:"domain"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Load users
		users, err := s.storage.LoadUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Find and update user
		for i, user := range users {
			if user.ID == userID {
				users[i].Username = userData.Username
				users[i].Login = userData.Login
				users[i].Domain = userData.Domain
				users[i].ModifiedAt = time.Now()

				// Handle password update (including migration and reuse of existing password)
				newEncryptedPassword, migrated, err := s.credManager.UpdateUserCredentialsWithMigration(
					users[i].EncryptedPassword, // Old encrypted password
					userData.Password,          // New plaintext password (can be empty)
				)
				if err != nil {
					http.Error(w, "Failed to update password: "+err.Error(), http.StatusInternalServerError)
					return
				}

				// Update the encrypted password (may be unchanged, migrated, or new)
				users[i].EncryptedPassword = newEncryptedPassword

				// If a new password was provided, update Windows Credential Store
				if userData.Password != "" {
					hosts, _ := s.storage.LoadHosts()
					for _, host := range hosts {
						if host.UserID == userID {
							s.credManager.StoreCredential(host.Address, userData.Username, userData.Password)
						}
					}
				}

				if migrated {
					logging.Log(debug, "User password was migrated from AES to DPAPI during edit")
				}

				// Save users
				if err := s.storage.SaveUsers(users); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				json.NewEncoder(w).Encode(users[i])
				return
			}
		}

		http.Error(w, "User not found", http.StatusNotFound)

	case "DELETE":
		// Load users
		users, err := s.storage.LoadUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Find and remove user
		for i, user := range users {
			if user.ID == userID {
				users = append(users[:i], users[i+1:]...)

				// Save users
				if err := s.storage.SaveUsers(users); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		http.Error(w, "User not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleHosts(w http.ResponseWriter, r *http.Request) {
	debug := true
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		hosts, err := s.storage.LoadHosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logging.Log(debug, "API GET: Loaded", len(hosts), "hosts from storage")
		if len(hosts) > 0 {
			// Debug: Show all hosts with their WindowWidth/WindowHeight values
			for i, host := range hosts {
				logging.Log(debug, "API GET: Host #", i, host.Name, "WindowWidth:", host.WindowWidth, "WindowHeight:", host.WindowHeight)
			}
		}
		json.NewEncoder(w).Encode(hosts)

	case "POST":
		var hostData struct {
			Address           string `json:"address"`
			Port              int    `json:"port"`
			UserID            string `json:"user_id"`
			DesktopWidth      int    `json:"desktop_width"`
			DesktopHeight     int    `json:"desktop_height"`
			WindowWidth       int    `json:"window_width"`
			WindowHeight      int    `json:"window_height"`
			PositionX         int    `json:"position_x"`
			PositionY         int    `json:"position_y"`
			RedirectClipboard bool   `json:"redirect_clipboard"`
			RedirectDrives    bool   `json:"redirect_drives"`
			DisplayMode       string `json:"display_mode"`
		}

		if err := json.NewDecoder(r.Body).Decode(&hostData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create new host
		host := models.NewHost(hostData.Address, hostData.Address, hostData.Port, hostData.UserID)
		host.DesktopWidth = hostData.DesktopWidth
		host.DesktopHeight = hostData.DesktopHeight
		host.WindowWidth = hostData.WindowWidth
		host.WindowHeight = hostData.WindowHeight
		host.PositionX = hostData.PositionX
		host.PositionY = hostData.PositionY
		host.RedirectClipboard = hostData.RedirectClipboard
		host.RedirectDrives = hostData.RedirectDrives
		host.DisplayMode = hostData.DisplayMode

		logging.Log(debug, "API POST: Received hostData - WindowWidth:", hostData.WindowWidth, "WindowHeight:", hostData.WindowHeight)
		logging.Log(debug, "API POST: Created host - WindowWidth:", host.WindowWidth, "WindowHeight:", host.WindowHeight)

		// Save host
		hosts, _ := s.storage.LoadHosts()
		hosts = append(hosts, host)
		if err := s.storage.SaveHosts(hosts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(host)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleHostByID(w http.ResponseWriter, r *http.Request) {
	debug := true
	w.Header().Set("Content-Type", "application/json")

	// Extract host ID from URL
	hostID := filepath.Base(r.URL.Path)

	switch r.Method {
	case "PUT":
		var hostData struct {
			Address           string `json:"address"`
			Port              int    `json:"port"`
			UserID            string `json:"user_id"`
			DesktopWidth      int    `json:"desktop_width"`
			DesktopHeight     int    `json:"desktop_height"`
			WindowWidth       int    `json:"window_width"`
			WindowHeight      int    `json:"window_height"`
			PositionX         int    `json:"position_x"`
			PositionY         int    `json:"position_y"`
			RedirectClipboard bool   `json:"redirect_clipboard"`
			RedirectDrives    bool   `json:"redirect_drives"`
			DisplayMode       string `json:"display_mode"`
		}

		if err := json.NewDecoder(r.Body).Decode(&hostData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Load hosts
		hosts, err := s.storage.LoadHosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Find and update host
		for i, host := range hosts {
			if host.ID == hostID {
				hosts[i].Address = hostData.Address
				hosts[i].Name = hostData.Address
				hosts[i].Port = hostData.Port
				hosts[i].UserID = hostData.UserID
				hosts[i].DesktopWidth = hostData.DesktopWidth
				hosts[i].DesktopHeight = hostData.DesktopHeight
				hosts[i].WindowWidth = hostData.WindowWidth
				hosts[i].WindowHeight = hostData.WindowHeight
				hosts[i].PositionX = hostData.PositionX
				hosts[i].PositionY = hostData.PositionY
				hosts[i].RedirectClipboard = hostData.RedirectClipboard
				hosts[i].RedirectDrives = hostData.RedirectDrives
				hosts[i].DisplayMode = hostData.DisplayMode
				hosts[i].ModifiedAt = time.Now()

				logging.Log(debug, "API PUT: Received hostData - WindowWidth:", hostData.WindowWidth, "WindowHeight:", hostData.WindowHeight)
				logging.Log(debug, "API PUT: Updated host - WindowWidth:", hosts[i].WindowWidth, "WindowHeight:", hosts[i].WindowHeight)

				// Save hosts
				if err := s.storage.SaveHosts(hosts); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				json.NewEncoder(w).Encode(hosts[i])
				return
			}
		}

		http.Error(w, "Host not found", http.StatusNotFound)

	case "DELETE":
		// Load hosts
		hosts, err := s.storage.LoadHosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Find and remove host
		for i, host := range hosts {
			if host.ID == hostID {
				hosts = append(hosts[:i], hosts[i+1:]...)

				// Save hosts
				if err := s.storage.SaveHosts(hosts); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		http.Error(w, "Host not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	debug := true
	logging.Log(debug, "API: Launching RDP connection")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var launchData struct {
		HostID string `json:"host_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&launchData); err != nil {
		logging.Log(true, "ERROR: Failed to decode launch data:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logging.Log(debug, "Launch request for host ID:", launchData.HostID)

	// Load data
	hosts, err := s.storage.LoadHosts()
	if err != nil {
		http.Error(w, "Failed to load hosts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := s.storage.LoadUsers()
	if err != nil {
		http.Error(w, "Failed to load users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Find host
	logging.Log(debug, "Searching for host in", len(hosts), "loaded hosts")
	var selectedHost *models.Host
	for _, host := range hosts {
		if host.ID == launchData.HostID {
			selectedHost = &host
			logging.Log(debug, "Found host:", host.Name, "at", host.Address)
			break
		}
	}

	if selectedHost == nil {
		logging.Log(true, "ERROR: Host not found with ID:", launchData.HostID)
		http.Error(w, "Host not found", http.StatusNotFound)
		return
	}

	// Find user
	logging.Log(debug, "Searching for user with ID:", selectedHost.UserID)
	var selectedUser *models.User
	for _, user := range users {
		if user.ID == selectedHost.UserID {
			selectedUser = &user
			logging.Log(debug, "Found user:", user.Username)
			break
		}
	}

	if selectedUser == nil {
		logging.Log(true, "ERROR: User not found with ID:", selectedHost.UserID)
		http.Error(w, "User not found for host", http.StatusNotFound)
		return
	}

	// Launch RDP
	logging.Log(debug, "Launching RDP connection for host:", selectedHost.Name, "user:", selectedUser.Username)
	if err := s.rdpGenerator.LaunchHost(*selectedHost, *selectedUser); err != nil {
		logging.Log(true, "ERROR: Failed to launch RDP:", err)
		http.Error(w, "Failed to launch RDP: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Log(debug, "RDP connection launched successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "launched"})
}

// saveUserAfterMigration callback for password migration
func (s *Server) saveUserAfterMigration(user models.User) error {
	users, err := s.storage.LoadUsers()
	if err != nil {
		return err
	}

	// Find and update user
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
			break
		}
	}

	return s.storage.SaveUsers(users)
}

// handleWindowInfo provides window border information for calculations
func (s *Server) handleWindowInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Get current window information by finding LaunchRDP window
		hwnd := findWindowByTitle("LaunchRDP")
		if hwnd == 0 {
			http.Error(w, "LaunchRDP window not found", http.StatusNotFound)
			return
		}

		borderInfo, err := GetWindowBorderInfo(hwnd)
		if err != nil {
			http.Error(w, "Failed to get window info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(borderInfo)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWebViewLog receives log messages from JavaScript and forwards them to backend logging
func (s *Server) handleWebViewLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logData struct {
		Message string `json:"message"`
		Level   string `json:"level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Forward to backend logging system
	debug := (logData.Level == "debug")
	logging.Log(debug, "[JavaScript]", logData.Message)

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
}
