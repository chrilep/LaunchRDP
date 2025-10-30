# 🚀 LaunchRDP# 🚀 LaunchRDP# 🚀 LaunchRDP# 🚀 LaunchRDP

[![Go Version](https://img.shields.io/github/go-mod/go-version/chrilep/LaunchRDP)](https://golang.org)[![Go Version](https://img.shields.io/github/go-mod/go-version/chrilep/LaunchRDP)](https://golang.org)

[![License](https://img.shields.io/github/license/chrilep/LaunchRDP)](https://github.com/chrilep/LaunchRDP/blob/main/LICENSE)[![License](https://img.shields.io/github/license/chrilep/LaunchRDP)](LICENSE)

[![Latest Release](https://img.shields.io/github/v/release/chrilep/LaunchRDP)](https://github.com/chrilep/LaunchRDP/releases)[![Release](https://img.shields.io/github/v/release/chrilep/LaunchRDP)](../../releases)

[![Open Issues](https://img.shields.io/github/issues/chrilep/LaunchRDP)](https://github.com/chrilep/LaunchRDP/issues)[![Issues](https://img.shields.io/github/issues/chrilep/LaunchRDP)](../../issues)

[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](https://github.com/chrilep/LaunchRDP/stargazers)[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and a **WebView2** frontend. It provides a clean, web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with precise window positioning, and launches connections with a single click.**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**.

## ✨ Core FeaturesIt provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.

- 🖥️ **Modern Web Interface**: A clean, responsive 4-column layout powered by WebView2.## ✨ Features

- 👥 **User Management**: Securely store and manage multiple user credentials.

- 🏠 **Host Management**: Configure RDP connections with a comprehensive set of options.**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**. **LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**.

- 📐 **Smart Window Positioning**: Automatically detects window borders via the Windows API to calculate the exact RDP client size and position.

- 🚀 **One-Click Launch**: Instantly generates and launches `.rdp` files using the native Microsoft RDP client (`mstsc.exe`).- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

- 🔒 **Local Storage**: Stores all configuration in simple JSON files within your user profile (`%APPDATA%`).

- ⚡ **Self-Contained**: A single, portable executable with no external dependencies required.- 👥 **User Management** - Securely store multiple user credentials It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.

## 🛠️ Installation- 🏠 **Host Management** - Configure RDP connections with full settings support

### Prerequisites- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioningIt provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.

- **Windows 10/11**

- **WebView2 Runtime**: This is included in modern Windows versions. If not, it can be [downloaded from Microsoft](https://developer.microsoft.com/en-us/microsoft-edge/webview2/).- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

### Download- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`## ✨ Features

You can download the latest pre-compiled binary from the [**Releases**](https://github.com/chrilep/LaunchRDP/releases) page.

- 📱 **Responsive Design** - Works seamlessly on different screen sizes

### Build from Source

If you prefer to build it yourself:- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

````bash

# 1. Clone the repository## 🔧 Technical Features- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

git clone https://github.com/chrilep/LaunchRDP.git

cd LaunchRDP- **Precise Window Calculations** - Uses Windows API to detect actual window borders- 👥 **User Management** - Securely store multiple user credentials ## ✨ Features



# 2. Build the application- **Complete RDP Template** - Supports all Microsoft RDP settings and options

go build -o LaunchRDP.exe .

```- **WinPosStr Generation** - Correctly calculates window positioning strings- 🏠 **Host Management** - Configure RDP connections with full settings support



## 🚦 Usage- **Real-time Updates** - Live calculation display as you adjust window dimensions



1.  **Launch**: Run `LaunchRDP.exe`. The application window will open automatically.- **Debug Logging** - Comprehensive logging for troubleshooting- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning**LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**. **LaunchRDP** is a modern Windows Remote Desktop Connection Manager built with **Go** and **WebView2**.

2.  **Add Users**: Navigate to the **Users** tab to add username/password credentials.

3.  **Add Hosts**: Go to the **Hosts** tab to create and configure your RDP connections.- **Version Management** - Integrated build system with automatic versioning

4.  **Connect**: Select a host from the list and click the launch button to start the session.

- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

### Interface Overview

- **📝 Edit Host (Column 1)**: Configure all RDP settings, including address, display mode, and window size.## 🛠️ Installation

- **🏠 Hosts (Column 2)**: The main list of all your saved RDP hosts.

- **👥 Users (Column 3)**: A list of all saved user credentials.- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

- **📊 Calculation Info**: A real-time display in the edit form showing the calculated RDP client size and `winposstr` value.

### Download

### Command Line Options

```bashDownload the latest release from the [Releases](../../releases) page.- 📱 **Responsive Design** - Works seamlessly on different screen sizes

# Start the application with default settings

.\LaunchRDP.exe### Requirements- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2- 👥 **User Management** - Securely store multiple user credentials It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.It provides a clean web-based interface to manage RDP hosts and users, automatically generates `.rdp` files with proper window positioning, and launches connections with a single click.



# Start on a custom port (default is 8088)- **Windows 10/11** (WebView2 runtime - usually pre-installed)

.\LaunchRDP.exe -port 9000

- **Microsoft Remote Desktop Connection** (`mstsc.exe`)## 🔧 Technical Features- 🏠 **Host Management** - Configure RDP connections with full settings support

# Show version information

.\LaunchRDP.exe -version### Build from Source- **Precise Window Calculations** - Uses Windows API to detect actual window borders- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning## ✨ Features## ✨ Features

````

```bash

### Storage Location

All host and user data is stored locally in your user's AppData directory:git clone https://github.com/chrilep/LaunchRDP.git- **Complete RDP Template** - Supports all Microsoft RDP settings and options

```

%APPDATA%\Lancer\LaunchRDP\cd LaunchRDP

├── hosts.json

└── users.jsongo build -o LaunchRDP.exe .- **WinPosStr Generation** - Correctly calculates window positioning strings- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

```

```

## 🏗️ Architecture

- **Backend**: A lightweight **Go** HTTP server that provides a JSON API.- **Real-time Updates** - Live calculation display as you adjust window dimensions

- **Frontend**: A modern UI built with standard **HTML, CSS, and JavaScript**.

- **UI Rendering**: **WebView2** is used to host the web-based frontend in a native window.## 🚦 Usage

- **RDP Integration**: Dynamically generates `.rdp` files and launches them with `mstsc.exe`.

- **System Integration**: Uses native Windows API calls for accurate window border detection.- **Debug Logging** - Comprehensive logging for troubleshooting- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2- 🖥️ **Modern Web Interface** - Clean, responsive 4-column layout using WebView2

## 🤝 Contributing### Quick Start

Contributions are welcome! Feel free to open an issue or submit a pull request.

- **Version Management** - Integrated build system with automatic versioning

1.  Fork the repository.

2.  Create your feature branch (`git checkout -b feature/MyNewFeature`).1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

3.  Commit your changes (`git commit -m 'Add some feature'`).

4.  Push to the branch (`git push origin feature/MyNewFeature`).2. **Add Users**: Click "Users" tab → Add username/password credentials- 📱 **Responsive Design** - Works seamlessly on different screen sizes

5.  Open a Pull Request.

6.  **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/chrilep/LaunchRDP/blob/main/LICENSE) file for details.4. **Connect**: Select a host from the list and click the connection button## 🛠️ Installation

### Interface Overview- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2- 👥 **User Management** - Securely store multiple user credentials - 👥 **User Management** - Securely store multiple user credentials

- **📝 Edit Host** (Column 1): Configure connection details, window size, and features### Download

- **🏠 Hosts** (Column 2): List of all configured RDP hosts

- **👥 Users** (Column 3): Manage user credentialsDownload the latest release from the [Releases](../../releases) page.## 🔧 Technical Features- 🏠 **Host Management** - Configure RDP connections with full settings support- 🏠 **Host Management** - Configure RDP connections with full settings support

- **📊 Calculation Info**: Real-time display of RDP client size and positioning

### Requirements- **Precise Window Calculations** - Uses Windows API to detect actual window borders- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning- 📐 **Smart Window Positioning** - Automatic window border detection and precise positioning

### Command Line Options

- **Windows 10/11** (WebView2 runtime - usually pre-installed)

```bash

LaunchRDP.exe              # Start with default settings- **Microsoft Remote Desktop Connection** (`mstsc.exe`)- **Complete RDP Template** - Supports all Microsoft RDP settings and options

LaunchRDP.exe -version     # Show version information

LaunchRDP.exe -port 8080   # Use custom port (default: 8088)### Build from Source- **WinPosStr Generation** - Correctly calculates window positioning strings- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly- 🚀 **One-Click Launch** - Generate `.rdp` files and launch connections instantly

```

````bash

### Storage Location

git clone https://github.com/chrilep/LaunchRDP.git- **Real-time Updates** - Live calculation display as you adjust window dimensions

Configuration files are stored in:

```cd LaunchRDP

%APPDATA%\Lancer\LaunchRDP\

├── hosts.json    # Host configurations  go build -o LaunchRDP.exe .- **Debug Logging** - Comprehensive logging for troubleshooting- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`- 🔒 **Secure Storage** - JSON-based local storage in `%APPDATA%\Lancer\LaunchRDP`

└── users.json    # User credentials

````

## 🎯 Key Features Implemented- **Version Management** - Integrated build system with automatic versioning

- ✅ **Modern WebView2 Interface** - Clean 4-column responsive layout## 🚦 Usage

- ✅ **Complete Host Management** - Full RDP configuration with all settings

- ✅ **User Credential Storage** - JSON-based local storage system - 📱 **Responsive Design** - Works seamlessly on different screen sizes- 📱 **Responsive Design** - Works seamlessly on different screen sizes

- ✅ **Smart Window Positioning** - Automatic border detection and calculations

- ✅ **Real-time Calculations** - Live display of client size and positioning### Quick Start

- ✅ **One-Click RDP Launch** - Generate and execute `.rdp` files instantly

- ✅ **Debug & Logging** - Comprehensive logging for troubleshooting## 🛠️ Installation

- ✅ **Build System** - Automated versioning and Windows resource embedding

1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

## 🔮 Future Enhancements

2. **Add Users**: Click "Users" tab → Add username/password credentials- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2- ⚡ **Zero Dependencies** - Self-contained executable with embedded WebView2

- 📁 **Connection Grouping** - Organize hosts by environment/purpose

- 📈 **Connection History** - Track recent connections and usage statistics 3. **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

- 💾 **Import/Export** - Backup and restore configurations

- 🔍 **Search & Filter** - Find hosts and users quickly4. **Connect**: Select a host from the list and click the connection button### Download

- 🎨 **Themes** - Dark/light mode and customization options

### Interface OverviewDownload the latest release from the [Releases](../../releases) page.## 🔧 Technical Features## 🔧 Technical Features

## 🏗️ Architecture

- **📝 Edit Host** (Column 1): Configure connection details, window size, and features### Requirements- **Precise Window Calculations** - Uses Windows API to detect actual window borders- **Precise Window Calculations** - Uses Windows API to detect actual window borders

LaunchRDP uses a clean separation of concerns:

- **🏠 Hosts** (Column 2): List of all configured RDP hosts

- **Frontend**: Modern HTML5/CSS3/JavaScript web interface

- **Backend**: Go HTTP server with JSON APIs - **👥 Users** (Column 3): Manage user credentials- **Windows 10/11** (WebView2 runtime - usually pre-installed)

- **Storage**: Local JSON files in `%APPDATA%`

- **RDP Generation**: Dynamic `.rdp` file creation with Microsoft-compliant templates- **📊 Calculation Info**: Real-time display of RDP client size and positioning

- **Window Management**: Native Windows API integration for border detection

- **Microsoft Remote Desktop Connection** (`mstsc.exe`)- **Complete RDP Template** - Supports all Microsoft RDP settings and options- **Complete RDP Template** - Supports all Microsoft RDP settings and options

## 📄 License

### Command Line Options

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

### Build from Source- **WinPosStr Generation** - Correctly calculates window positioning strings- **WinPosStr Generation** - Correctly calculates window positioning strings

## 🤝 Contributing

````bash

Contributions, issues, and feature requests are welcome!

LaunchRDP.exe              # Start with default settings```bash

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/amazing-feature`)LaunchRDP.exe -version     # Show version information

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)LaunchRDP.exe -port 8080   # Use custom port (default: 8088)git clone https://github.com/chrilep/LaunchRDP.git- **Real-time Updates** - Live calculation display as you adjust window dimensions- **Real-time Updates** - Live calculation display as you adjust window dimensions

5. Open a Pull Request

````

## 🙏 Acknowledgements

cd LaunchRDP

- **Go** - Programming language and standard library

- **WebView2** - Modern web runtime for Windows### Storage Location

- **Microsoft RDP** - Remote Desktop Protocol and `mstsc.exe`

- **Windows API** - Native system integrationgo build -o LaunchRDP.exe .- **Debug Logging** - Comprehensive logging for troubleshooting- **Debug Logging** - Comprehensive logging for troubleshooting

Configuration files are stored in:

``````

%APPDATA%\Lancer\LaunchRDP\

├── hosts.json    # Host configurations  - **Version Management** - Integrated build system with automatic versioning- **Version Management** - Integrated build system with automatic versioning

└── users.json    # User credentials

```## 🚦 Usage



## 🎯 Key Features Implemented## 🛠️ Installation## 🛠️ Installation



- ✅ **Modern WebView2 Interface** - Clean 4-column responsive layout### Quick Start

- ✅ **Complete Host Management** - Full RDP configuration with all settings

- ✅ **User Credential Storage** - JSON-based local storage system  ### Download### Download

- ✅ **Smart Window Positioning** - Automatic border detection and calculations

- ✅ **Real-time Calculations** - Live display of client size and positioning1. **Launch**: Double-click `LaunchRDP.exe` (opens automatically in WebView2)

- ✅ **One-Click RDP Launch** - Generate and execute `.rdp` files instantly

- ✅ **Debug & Logging** - Comprehensive logging for troubleshooting2. **Add Users**: Click "Users" tab → Add username/password credentialsDownload the latest release from the [Releases](../../releases) page.

- ✅ **Build System** - Automated versioning and Windows resource embedding

3. **Add Hosts**: Click "Hosts" tab → Configure RDP connection settings

## 🔮 Future Enhancements

4. **Connect**: Select a host from the list and click the connection buttonDownload the latest release from the [Releases](../../releases) page.

- 📁 **Connection Grouping** - Organize hosts by environment/purpose

- 📈 **Connection History** - Track recent connections and usage statistics  ### Interface Overview### Requirements

- 💾 **Import/Export** - Backup and restore configurations

- 🔍 **Search & Filter** - Find hosts and users quickly- **📝 Edit Host** (Column 1): Configure connection details, window size, and features- **Windows 10/11** (WebView2 runtime - usually pre-installed)### Requirements

- 🎨 **Themes** - Dark/light mode and customization options

- **🏠 Hosts** (Column 2): List of all configured RDP hosts

## 🏗️ Architecture

- **👥 Users** (Column 3): Manage user credentials- **Microsoft Remote Desktop Connection** (`mstsc.exe`)

LaunchRDP uses a clean separation of concerns:

- **📊 Calculation Info**: Real-time display of RDP client size and positioning

- **Frontend**: Modern HTML5/CSS3/JavaScript web interface

- **Backend**: Go HTTP server with JSON APIs  - **Windows 10/11** (WebView2 runtime - usually pre-installed)

- **Storage**: Local JSON files in `%APPDATA%`

- **RDP Generation**: Dynamic `.rdp` file creation with Microsoft-compliant templates### Command Line Options

- **Window Management**: Native Windows API integration for border detection

### Build from Source- **Microsoft Remote Desktop Connection** (`mstsc.exe`)

## 📄 License

`````bash

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

LaunchRDP.exe              # Start with default settings````bash

## 🤝 Contributing

LaunchRDP.exe -version     # Show version information

Contributions, issues, and feature requests are welcome!

LaunchRDP.exe -port 8080   # Use custom port (default: 8088)git clone https://github.com/chrilep/LaunchRDP.git### Build from Source

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/amazing-feature`)`````

3. Commit your changes (`git commit -m 'Add amazing feature'`)

4. Push to the branch (`git push origin feature/amazing-feature`)cd LaunchRDP

5. Open a Pull Request

### Storage Location

## 🙏 Acknowledgements

go build -o LaunchRDP.exe .```bash

- **Go** - Programming language and standard library

- **WebView2** - Modern web runtime for WindowsConfiguration files are stored in:

- **Microsoft RDP** - Remote Desktop Protocol and `mstsc.exe`

- **Windows API** - Native system integration`````git clone https://github.com/chrilep/LaunchRDP.git

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
``````
