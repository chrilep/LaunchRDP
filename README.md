# 🚀 LaunchRDP

[![Go Version](https://img.shields.io/github/go-mod/go-version/chrilep/LaunchRDP)](https://golang.org)
[![License](https://img.shields.io/github/license/chrilep/LaunchRDP)](LICENSE)
[![Release](https://img.shields.io/github/v/release/chrilep/LaunchRDP)](../../releases)
[![Issues](https://img.shields.io/github/issues/chrilep/LaunchRDP)](../../issues)
[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

**LaunchRDP** is a Windows utility written in **Go** with a **Web UI**.  
It manages credentials (username + password) and host definitions, then dynamically generates `.rdp` files with all supported settings. With one click, LaunchRDP launches sessions directly using **mstsc.exe** (Microsoft Remote Desktop Connection).

## 🔒 Security-First Design

- **Zero 3rd party dependencies** - Uses only Go standard library
- **Minimal attack surface** - No external GUI frameworks or libraries
- **Native Windows integration** - Direct Windows API calls for credentials
- **Local-only web interface** - No network exposure by default

## ⚠️ Project Status

> **Note:** LaunchRDP is in the **early stages of development**.  
> Features, UI, and storage formats are subject to change. Expect breaking changes until the first stable release.

## ✨ Planned Features

- Manage multiple **credentials** (user + password) securely in Windows Credential Store
- Define **hosts** with full `.rdp` configuration options
- Dynamically generate temporary `.rdp` files on demand
- One‑click session launch via `mstsc.exe`
- Modern responsive **web interface** - works on any device
- **Zero external dependencies** - pure Go standard library
- Automatic browser launch or manual access via `http://localhost:8888`
- Cross-platform potential (currently Windows-focused)

## 🛠️ Installation

Precompiled binaries will be published on the [Releases](../../releases) page once available.

## 🚦 Usage

### Quick Start

1. **Run the application**: `LaunchRDP.exe` (opens web interface automatically)
2. **Add users**: Go to Users tab, add credentials (username + password)
3. **Add hosts**: Go to Hosts tab, configure RDP connections
4. **Launch connections**: Go to Launch tab, click any connection button

### Command Line Options

```bash
LaunchRDP.exe -port 8888    # Custom port (default: 8080)
LaunchRDP.exe -version      # Show version information
```

### Web Interface

- **Modern UI**: Clean, responsive design that works on desktop and mobile
- **Real-time updates**: Changes are saved automatically
- **Secure storage**: Passwords encrypted with Windows DPAPI
- **No installation**: Self-contained executable

## 📌 Roadmap

- [x] **Credential storage** - Windows Credential Manager integration
- [x] **Web UI** - Modern responsive interface for managing hosts and credentials
- [x] **RDP settings** - Screen size, clipboard redirection, and more
- [x] **Auto-save** - Changes saved automatically without manual save buttons
- [x] **Zero dependencies** - Pure Go standard library implementation
- [ ] **Profiles and grouping** - Organize connections by environment/purpose
- [ ] **Connection history** - Track recent connections and favorites
- [ ] **Import/Export** - Backup and restore configurations

## 📄 License

This project is licensed under the MIT License — see the LICENSE file for details.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!
Feel free to open an issue or submit a pull request.

## 🙏 Acknowledgements

- Go
- Fyne
- Microsoft’s mstsc.exe client
