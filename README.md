# 🚀 LaunchRDP

[![Go Version](https://img.shields.io/github/go-mod/go-version/chrilep/LaunchRDP)](https://golang.org)
[![License](https://img.shields.io/github/license/chrilep/LaunchRDP)](LICENSE)
[![Release](https://img.shields.io/github/v/release/chrilep/LaunchRDP)](../../releases)
[![Issues](https://img.shields.io/github/issues/chrilep/LaunchRDP)](../../issues)
[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**.

It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.

[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

## ✨ Features

**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**. **LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**.

- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

- 👥 **User Management** - Securely store multiple user credentials It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.

- 🏠 **Host Management** - Configure RDP connections with full settings support

- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning## ✨ Features## ✨ Features

- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

- 📱 **Responsive Design** - Works seamlessly on different screen sizes

- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2- 👥 **User Management** - Securely store multiple user credentials - 👥 **User Management** - Securely store multiple user credentials

## 🔧 Technical Features- 🏠 **Host Management** - Configure RDP connections with full settings support- 🏠 **Host Management** - Configure RDP connections with full settings support

- **Precise Window Calculations** - Uses Windows API to detect actual window borders- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning

- **Complete RDP Template** - Supports all Microsoft RDP settings and options

- **WinPosStr Generation** - Correctly calculates window positioning strings- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

- **Real-time Updates** - Live calculation display as you adjust window dimensions

- **Debug Logging** - Comprehensive logging for troubleshooting- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`

- **Version Management** - Integrated build system with automatic versioning

- 📱 **Responsive Design** - Works seamlessly on different screen sizes- 📱 **Responsive Design** - Works seamlessly on different screen sizes

## 🛠️ Installation

- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2

### Download

Download the latest release from the [Releases](../../releases) page.## 🔧 Technical Features## 🔧 Technical Features

### Requirements- **Precise Window Calculations** - Uses Windows API to detect actual window borders- **Precise Window Calculations** - Uses Windows API to detect actual window borders

- **Windows 10/11** (WebView2 runtime - usually pre-installed)

- **Microsoft Remote Desktop Connection** (`mstsc.exe`)- **Complete RDP Template** - Supports all Microsoft RDP settings and options- **Complete RDP Template** - Supports all Microsoft RDP settings and options

### Build from Source- **WinPosStr Generation** - Correctly calculates window positioning strings- **WinPosStr Generation** - Correctly calculates window positioning strings

```bash

git clone https://github.com/chrilep/LaunchRDP.git- **Real-time Updates** - Live calculation display as you adjust window dimensions- **Real-time Updates** - Live calculation display as you adjust window dimensions

cd LaunchRDP

go build -o LaunchRDP.exe .- **Debug Logging** - Comprehensive logging for troubleshooting- **Debug Logging** - Comprehensive logging for troubleshooting

```

- **Version Management** - Integrated build system with automatic versioning- **Version Management** - Integrated build system with automatic versioning

## 🚦 Usage

## 🛠️ Installation## 🛠️ Installation

### Quick Start

### Download### Download

1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

2. **Add Users**: Click "Users" tab → Add username/password credentialsDownload the latest release from the [Releases](../../releases) page.

3. **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

4. **Connect**: Select a host from the list and click the connection buttonDownload the latest release from the [Releases](../../releases) page.

### Interface Overview### Requirements

- **📝 Edit Host** (Column 1): Configure connection details, window size, and features- **Windows 10/11** (WebView2 runtime - usually pre-installed)### Requirements

- **🏠 Hosts** (Column 2): List of all configured RDP hosts

- **👥 Users** (Column 3): Manage user credentials- **Microsoft Remote Desktop Connection** (`mstsc.exe`)

- **📊 Calculation Info**: Real-time display of RDP client size and positioning

- **Windows 10/11** (WebView2 runtime - usually pre-installed)

### Command Line Options

### Build from Source- **Microsoft Remote Desktop Connection** (`mstsc.exe`)

`````bash

LaunchRDP.exe              # Start with default settings````bash

LaunchRDP.exe -version     # Show version information

LaunchRDP.exe -port 8080   # Use custom port (default: 8088)git clone https://github.com/chrilep/LaunchRDP.git### Build from Source

`````

cd LaunchRDP

### Storage Location

go build -o LaunchRDP.exe .```bash

Configuration files are stored in:

`````git clone https://github.com/chrilep/LaunchRDP.git

%APPDATA%\Lancer\LaunchRDP\

├── hosts.json    # Host configurations  cd LaunchRDP

└── users.json    # User credentials

```## 🚦 Usagego build -o LaunchRDP.exe .



## 🎯 Key Features Implemented````



- ✅ **Modern WebView2 Interface** - Clean 4-column responsive layout### Quick Start

- ✅ **Complete Host Management** - Full RDP configuration with all settings

- ✅ **User Credential Storage** - JSON-based local storage system  ## 🚦 Usage

- ✅ **Smart Window Positioning** - Automatic border detection and calculations

- ✅ **Real-time Calculations** - Live display of client size and positioning1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

- ✅ **One-Click RDP Launch** - Generate and execute `.rdp` files instantly

- ✅ **Debug & Logging** - Comprehensive logging for troubleshooting2. **Add Users**: Click "Users" tab → Add username/password credentials### Quick Start

- ✅ **Build System** - Automated versioning and Windows resource embedding

3. **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

## 🔮 Future Enhancements

4. **Connect**: Select a host from the list and click the connection button1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

- 📁 **Connection Grouping** - Organize hosts by environment/purpose

- 📈 **Connection History** - Track recent connections and usage statistics  5. **Add Users**: Click "Users" tab → Add username/password credentials

- 💾 **Import/Export** - Backup and restore configurations

- 🔍 **Search & Filter** - Find hosts and users quickly### Interface Overview3. **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

- 🎨 **Themes** - Dark/light mode and customization options

4. **Connect**: Select a host from the list and click the connection button

## 🏗️ Architecture

- **📝 Edit Host** (Column 1): Configure connection details, window size, and features

LaunchRDP uses a clean separation of concerns:

- **🏠 Hosts** (Column 2): List of all configured RDP hosts ### Interface Overview

- **Frontend**: Modern HTML5/CSS3/JavaScript web interface

- **Backend**: Go HTTP server with JSON APIs  - **👥 Users** (Column 3): Manage user credentials

- **Storage**: Local JSON files in `%APPDATA%`

- **RDP Generation**: Dynamic `.rdp` file creation with Microsoft-compliant templates- **📊 Calculation Info**: Real-time display of RDP client size and positioning- **📝 Edit Host** (Column 1): Configure connection details, window size, and features

- **Window Management**: Native Windows API integration for border detection

- **🏠 Hosts** (Column 2): List of all configured RDP hosts

## 📄 License

### Command Line Options- **👥 Users** (Column 3): Manage user credentials

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

- **📊 Calculation Info**: Real-time display of RDP client size and positioning

## 🤝 Contributing

````bash

Contributions, issues, and feature requests are welcome!

LaunchRDP.exe              # Start with default settings### Command Line Options

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/amazing-feature`)LaunchRDP.exe -version     # Show version information

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)LaunchRDP.exe -port 8080   # Use custom port (default: 8088)```bash

5. Open a Pull Request

```LaunchRDP.exe              # Start with default settings

## 🙏 Acknowledgements

LaunchRDP.exe -version     # Show version information

- **Go** - Programming language and standard library

- **WebView2** - Modern web runtime for Windows### Storage LocationLaunchRDP.exe -port 8080   # Use custom port (default: 8088)

- **Microsoft RDP** - Remote Desktop Protocol and `mstsc.exe`

- **Windows API** - Native system integration````

Configuration files are stored in:

```### Storage Location

%APPDATA%\Lancer\LaunchRDP\

├── hosts.json    # Host configurations  Configuration files are stored in:

└── users.json    # User credentials

```

%APPDATA%\Lancer\LaunchRDP\

## 🎯 Key Features Implemented├── hosts.json # Host configurations

└── users.json # User credentials

- ✅ **Modern WebView2 Interface** - Clean 4-column responsive layout```

- ✅ **Complete Host Management** - Full RDP configuration with all settings

- ✅ **User Credential Storage** - JSON-based local storage system ## 📌 Roadmap

- ✅ **Smart Window Positioning** - Automatic border detection and calculations

- ✅ **Real-time Calculations** - Live display of client size and positioning- [x] **Credential storage** - Windows Credential Manager integration

- ✅ **One-Click RDP Launch** - Generate and execute `.rdp` files instantly- [x] **Web UI** - Modern responsive interface for managing hosts and credentials

- ✅ **Debug & Logging** - Comprehensive logging for troubleshooting- [x] **RDP settings** - Screen size, clipboard redirection, and more

- ✅ **Build System** - Automated versioning and Windows resource embedding- [x] **Auto-save** - Changes saved automatically without manual save buttons

- [x] **Zero dependencies** - Pure Go standard library implementation

## 🔮 Future Enhancements- [ ] **Profiles and grouping** - Organize connections by environment/purpose

- [ ] **Connection history** - Track recent connections and favorites

- 📁 **Connection Grouping** - Organize hosts by environment/purpose- [ ] **Import/Export** - Backup and restore configurations

- 📈 **Connection History** - Track recent connections and usage statistics

- 💾 **Import/Export** - Backup and restore configurations## 📄 License

- 🔍 **Search & Filter** - Find hosts and users quickly

- 🎨 **Themes** - Dark/light mode and customization optionsThis project is licensed under the MIT License — see the LICENSE file for details.

## 🏗️ Architecture## 🤝 Contributing

LaunchRDP uses a clean separation of concerns:Contributions, issues, and feature requests are welcome!

Feel free to open an issue or submit a pull request.

- **Frontend**: Modern HTML5/CSS3/JavaScript web interface

- **Backend**: Go HTTP server with JSON APIs ## 🙏 Acknowledgements

- **Storage**: Local JSON files in `%APPDATA%`

- **RDP Generation**: Dynamic `.rdp` file creation with Microsoft-compliant templates- Go

- **Window Management**: Native Windows API integration for border detection- Fyne

- Microsoft’s mstsc.exe client

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgements

- **Go** - Programming language and standard library
- **WebView2** - Modern web runtime for Windows
- **Microsoft RDP** - Remote Desktop Protocol and `mstsc.exe`
- **Windows API** - Native system integration
`````
