# Changelog

All notable changes to LaunchRDP will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Enterprise Deployment Support**
  - Professional NSIS installer with silent installation support
  - Silent install via `/S` parameter for SCCM/Intune deployment
  - Optional desktop shortcut creation during installation
  - Multi-user data cleanup options in uninstaller
  - Group Policy (GPO) deployment compatible
  - WebView2 runtime auto-detection and installation

- **Build System Enhancements**
  - Automated version management with `version.go` as single source of truth
  - PowerShell build script (`build.ps1`) with version increment workflow
  - Automatic icon preparation and integration
  - NSIS installer generation via `wails build -nsis`

### Changed
- Improved build process reliability and error handling
- Enhanced installer user experience with single-page UI
- Optimized installer size (~10.3 MB)

## [2.0.1] - 2025-11-09

### Major Changes
- **Complete migration from Gin HTTP server to Wails v2 framework**
  - Native Windows desktop application with WebView2
  - Eliminated HTTP server overhead and security concerns
  - Direct Go-to-JavaScript bindings for improved performance

### Added
- **RDP Window Reuse System**
  - Automatically detects existing RDP connections by address
  - Brings existing windows to front instead of launching duplicates
  - Uses Windows API (EnumWindows, TscShellContainerClass) for window detection
  - Restores minimized RDP windows (SW_RESTORE) before bringing to foreground

- **Enhanced Window Management**
  - Fixed Windows event hook (WINEVENT_INCONTEXT flag) for geometry tracking
  - Proper window position and size change detection
  - Window state synchronization with UI

- **Native Windows Credential Manager Integration**
  - Replaced cmdkey.exe with native API (CredWriteW/CredDeleteW)
  - Secure credential storage using CRED_TYPE_DOMAIN_PASSWORD
  - Proper UTF-16 encoding without null terminators
  - Automatic credential updates when host address/user changes
  - DPAPI encryption for JSON backup storage only (no credential reading to avoid antivirus flags)

- **Complete English Localization**
  - All source code comments translated to English
  - All UI elements in English
  - Prepared for Wails i18n integration in future updates

- **Build System Improvements**
  - `version.go` as single source of truth for version management
  - PowerShell build script with automatic wails.json synchronization
  - Comprehensive BUILD.md documentation
  - Clean version increment workflow (auto-increment build number)

### Changed
- Removed all popup notifications for cleaner UI experience
  - Replaced with console logging for debugging
  - No more obstructive messages during window switching

- Optimized credential handling
  - Credentials loaded by Windows automatically during RDP launch
  - Removed unnecessary DPAPI decryption on connection launch
  - Host editing properly updates credentials in Credential Manager

### Removed
- **System User / Pass-through Authentication**
  - Removed "Current Windows User" functionality (incompatible with Microsoft Accounts)
  - All users now require explicit username + password storage
  - Simplified user management without special system user handling

- HTTP server and all REST API endpoints
- Port configuration (no longer needed without HTTP server)
- Server-side routing logic

### Fixed
- User selection in Edit Host dialog now shows correct assigned user
- Credential storage Error 87 (invalid parameters) - fixed hostname\username format
- RDP password authentication failing - removed null terminator from UTF-16 passwords
- Window hook not firing - switched to WINEVENT_INCONTEXT flag

### Technical Details
- **Framework**: Wails v2.10.2
- **Go Version**: 1.21+
- **Frontend**: Vanilla JavaScript + WebView2
- **Build Target**: Windows amd64
- **Dependencies**: 
  - Native Windows APIs (advapi32.dll, user32.dll)
  - Windows Credential Manager
  - DPAPI for encryption

### Developer Notes
- `.gitignore` updated for Wails + Go + Node.js development
- Comprehensive coverage of build artifacts, IDE files, and OS-specific files
- `frontend/wailsjs/` kept in version control (auto-generated but required)

---

## [1.4.0] - 2025-11-06

### Added
- Backend logging integration for JavaScript console messages
- Improved debugging capabilities with unified logging system

---

## [1.3.0] - 2025-11-05

### Changed
- Internal improvements and optimizations
- Code cleanup and refactoring

---

## [1.3.0] - 2025-11-05

### Added
- Configured start position and size for RDP windows
- Window geometry persistence between sessions
- Custom window positioning support

### Changed
- Enhanced window state management
- Improved position calculation for multi-monitor setups

---

## [1.2.0] - 2025-11-04

### Added
- Enhanced type exclusions in credential handling
- Improved error handling for credential operations

### Changed
- Refined credential storage logic
- Better type safety in credential management

---

## [1.1.0] - 2025-11-04

### Added
- Initial credential type exclusions
- Base credential management system

### Changed
- Core credential handling improvements
- Storage optimization

---

## [1.0.0] - 2025-10-30

### Initial Release

#### Features
- **RDP Connection Management**
  - Store and manage multiple RDP hosts
  - Custom names for easy identification
  - Port configuration per host
  - Display mode selection (fullscreen/window)

- **User Credential Management**
  - Secure credential storage
  - Multiple user profiles
  - Domain support
  - Password encryption

- **Window Configuration**
  - Custom window sizes per host
  - Position configuration
  - Multi-monitor support in fullscreen mode
  - Fullscreen and windowed modes

- **Advanced Settings**
  - Clipboard redirection
  - Drive mapping
  - Custom RDP port configuration
  - Connection-specific settings

- **User Interface**
  - Clean, modern design
  - Dark mode support (automatic based on system preference)
  - Responsive layout
  - Column-based navigation

#### Technology Stack
- Backend: Go with Gin HTTP framework
- Frontend: HTML/CSS/JavaScript
- Storage: JSON file-based storage
- Security: DPAPI encryption for sensitive data
- Platform: Windows 10+ (WebView2 required)

#### Architecture
- REST API backend on configurable port
- WebView2-based frontend
- Local file storage in user AppData
- Secure credential management with Windows integration

---

## Version History Summary

- **v2.0.1+** (Unreleased) - Enterprise deployment support, NSIS installer, silent install capability
- **v2.0.1** - Major Wails v2 migration, native desktop app, RDP window reuse, native credential API
- **v1.4.0** - Backend logging improvements
- **v1.3.0** - Window position and size configuration
- **v1.2.0** - Credential handling enhancements
- **v1.1.0** - Credential management improvements
- **v1.0.0** - Initial release with HTTP server

---

## Upgrade Notes

### From v1.x to v2.0.1+

**Breaking Changes:**
- Application now runs as native desktop app (no HTTP server)
- Port configuration removed (no longer applicable)
- System user / pass-through authentication removed
- All stored credentials remain compatible (DPAPI + Windows Credential Manager)

**Migration Steps:**
1. Backup your data folder (typically `%APPDATA%\Lancer\LaunchRDP\`)
2. Install v2.0.1+ using NSIS installer or portable executable
3. First launch will automatically migrate existing data
4. Verify all hosts and users are present
5. Test RDP connections

**New Features Available:**
- RDP window reuse (no duplicate connections)
- Native Windows credential integration
- Improved window management
- Cleaner UI without popups
- Enterprise deployment support with silent installation
- Professional NSIS installer package

---

## Future Roadmap

### Planned Features
- [ ] Multi-language support using Wails i18n
- [ ] Host/User Cloning

---

## Contributing

See [GitHub Issues](https://github.com/chrilep/LaunchRDP/issues) for known issues and feature requests.

## License

See [LICENSE](LICENSE) file for details.
