# LaunchRDP

<div align="center">

![LaunchRDP](res/icon.png)

**Modern RDP Connection Manager for Windows**

A fast, secure, and user-friendly Remote Desktop Protocol (RDP) connection manager built with Wails v2 and Go.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.0.0-green.svg)](CHANGELOG.md)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![Wails](https://img.shields.io/badge/wails-v2.10.2-blue.svg)](https://wails.io/)

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Building](#-building) • [Changelog](CHANGELOG.md)

</div>

---

## 🎯 Why LaunchRDP?

LaunchRDP is designed as a modern alternative to legacy RDP managers like mRemoteNG. It provides:

- ✅ **Native Windows Integration** - Built with Wails v2, no HTTP server overhead
- ✅ **Secure Credential Storage** - Windows Credential Manager + DPAPI encryption
- ✅ **Smart Window Management** - Automatic RDP window reuse, position persistence
- ✅ **Clean Modern UI** - Dark mode support, responsive design
- ✅ **Zero Configuration** - Works out of the box, no complex setup

**Perfect for:**
- IT Professionals managing multiple servers
- System administrators with frequent RDP connections
- Anyone tired of typing credentials repeatedly
- Users seeking a modern replacement for outdated RDP managers

---

## ✨ Features

### Connection Management
- 🖥️ **Multi-Host Support** - Store unlimited RDP connections
- 👤 **User Profiles** - Manage multiple credential sets
- 🔐 **Secure Credentials** - Native Windows Credential Manager integration
- 📝 **Custom Names** - Friendly aliases for easy identification
- 🔄 **Window Reuse** - Automatically detects and activates existing connections

### Window & Display
- 📐 **Custom Positioning** - Save window positions per connection
- 📏 **Size Presets** - Configure window dimensions for each host
- 🖼️ **Fullscreen Mode** - Toggle between windowed and fullscreen
- 🖥️ **Multi-Monitor** - Full multi-monitor support

### Advanced Options
- 📋 **Clipboard Sharing** - Seamless copy/paste between local and remote
- 💾 **Drive Mapping** - Share local drives with remote sessions
- 🔌 **Custom Ports** - Configure non-standard RDP ports
- ⚙️ **Per-Connection Settings** - Individual configuration for each host

### Security & Privacy
- 🔒 **DPAPI Encryption** - Windows Data Protection API for stored passwords
- 🛡️ **Native Credential Storage** - Leverages Windows Credential Manager
- 🚫 **No Cloud Sync** - All data stays on your local machine
- 🔐 **Domain Support** - Full support for domain credentials

---

## 📋 Requirements

- **OS**: Windows 10/11 (64-bit)
- **Runtime**: WebView2 (usually pre-installed on Windows 10+)
- **Disk Space**: ~10 MB

> **Note**: WebView2 is automatically installed on Windows 11 and recent Windows 10 updates. If needed, it will download automatically on first launch.

---

## 🚀 Installation

### Option 1: Download Release (Recommended)

1. Download the latest `LaunchRDP.exe` from [Releases](https://github.com/chrilep/LaunchRDP/releases)
2. Run the executable - no installation needed!
3. Your data will be stored in `%APPDATA%\Lancer\LaunchRDP\`

### Option 2: Build from Source

See [Building from Source](#-building-from-source) section below.

---

## 💡 Usage

### Quick Start

1. **Launch the Application**
   - Run `LaunchRDP.exe`
   - The application window will open

2. **Add Your First User**
   - Click the **Users** tab
   - Click **Add User**
   - Enter username, login, domain (optional), and password
   - Click **OK**

3. **Add Your First Host**
   - Click the **Hosts** tab
   - Click **Add Host**
   - Enter custom name (optional) and host address/IP
   - Select the user credentials to use
   - Configure window settings if needed
   - Click **OK**

4. **Connect**
   - Click the **Launch** button next to any host
   - LaunchRDP will automatically use stored credentials
   - Window position and size will be restored if previously saved

### Interface Overview

```
┌─────────────────────────────────────┐
│         LaunchRDP                   │  ← Header
├─────────────┬───────────────────────┤
│   Hosts   │   Users   │            │  ← Navigation Tabs
├───────────────────────────────────┬─┤
│                                   │E│
│  🖥️ Production Server            │d│  ← Host List
│     192.168.1.10:3389            │i│
│     [Edit] [Launch]              │t│
│                                   │ │
│  🖥️ Test Environment             │H│  ← Edit Panel
│     test.domain.local            │o│     (shows when editing)
│     [Edit] [Launch]              │s│
│                                   │t│
└───────────────────────────────────┴─┘
```

### Tips & Tricks

- **Quick Launch**: Double-click a host to launch immediately
- **Window Reuse**: LaunchRDP automatically detects existing RDP windows and brings them to front instead of creating duplicates
- **Edit Shortcuts**: Click a host name to quickly edit its settings
- **Domain Credentials**: Use `DOMAIN\username` format in the login field
- **Password Update**: Leave password empty when editing users to keep existing password

---

## 🛠️ Building from Source

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 16+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

```powershell
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Build Steps

```powershell
# Clone the repository
git clone https://github.com/chrilep/LaunchRDP.git
cd LaunchRDP

# Build with the provided script
.\build.ps1

# Or build manually
wails build
```

The executable will be created in `build/bin/LaunchRDP.exe`

### Development Mode

```powershell
# Run in development mode with hot reload
.\build.ps1 dev

# Or manually
wails dev
```

For detailed build instructions, see [BUILD.md](BUILD.md)

---

## 📁 Project Structure

```
LaunchRDP/
├── app/                    # Backend Go code
│   ├── config/            # Configuration management
│   ├── credentials/       # Windows Credential Manager integration
│   ├── logging/           # Logging utilities
│   ├── models/            # Data models
│   ├── rdp/               # RDP file generation and launching
│   └── storage/           # JSON file storage
├── frontend/              # Frontend code
│   ├── src/              # JavaScript source
│   ├── index.html        # Main UI
│   └── style.css         # Styling
├── build/                # Build outputs
├── res/                  # Resources (icons)
├── app.go               # Main application logic
├── main.go              # Entry point
├── version.go           # Version information
├── wails.json           # Wails configuration
└── build.ps1            # Build script
```

---

## 🔧 Configuration

LaunchRDP stores all data locally:

- **Application Data**: `%APPDATA%\Lancer\LaunchRDP\data\`
  - `hosts.json` - Host configurations
  - `users.json` - User credentials (DPAPI encrypted)
  - `window_state.json` - Window position and size

- **Credentials**: Windows Credential Manager
  - Target: `rdp:{hostname}`
  - Automatically managed by the application

- **Logs**: `%APPDATA%\Lancer\LaunchRDP\logs\`

---

## 🆙 Upgrading from v1.x

### What Changed in v2.0

- ✅ Migrated from HTTP server to native Wails v2 application
- ✅ Improved security with native Credential Manager API
- ✅ Added RDP window reuse functionality
- ✅ Removed all popup notifications for cleaner UX
- ✅ Complete English localization

### Migration Steps

1. Backup your data: `%APPDATA%\Lancer\LaunchRDP\`
2. Install LaunchRDP v2.0.0
3. Launch the application - data migrates automatically
4. Verify all hosts and credentials

See [CHANGELOG.md](CHANGELOG.md) for detailed changes.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Reporting Issues

Found a bug? Have a feature request? Please [open an issue](https://github.com/chrilep/LaunchRDP/issues).

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Wails](https://wails.io/) - Go + Web frontend framework
- Inspired by mRemoteNG and other RDP managers
- Thanks to the Go and Wails communities

---

## 📧 Support

- **Issues**: [GitHub Issues](https://github.com/chrilep/LaunchRDP/issues)
- **Discussions**: [GitHub Discussions](https://github.com/chrilep/LaunchRDP/discussions)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)

---

<div align="center">

**⭐ If you find LaunchRDP useful, please consider giving it a star! ⭐**

Made with ❤️ by [Lancer](https://github.com/chrilep)

</div>
