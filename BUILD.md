# LaunchRDP Build System

## Version Management

**Single Source of Truth: `version.go`**

All version information is managed through `version.go` as the leading source:

```go
const (
    AppName    = "LaunchRDP"
    ID         = "com.chrilep.launchrdp"
    Version    = "2.0.0.1"  // ← ONLY EDIT THIS
    Author     = "Lancer"
    Repository = "https://github.com/chrilep/LaunchRDP"
)
```

### Version Format

- Format: `MAJOR.MINOR.PATCH.BUILD`
- Example: `2.0.0.1`
  - `2` = Major version
  - `0` = Minor version  
  - `0` = Patch version
  - `1` = Build number

### Automatic Synchronization

The build script automatically syncs versions to:
- ✅ `wails.json` → `info.productVersion` (used for Windows manifest)
- ✅ `build/windows/info.json` → File version metadata (via Wails template)
- ✅ `build/windows/wails.exe.manifest` → Assembly version (via Wails template)

**DO NOT manually edit version numbers in `wails.json` - they will be overwritten!**

## Build Commands

### Build Production Release
```powershell
.\build.ps1
```
- Prompts for new version (or press Enter to auto-increment build number)
- Updates `version.go` and syncs `wails.json`
- Runs `go mod tidy` and `wails generate`
- Builds production executable with `-clean` flag

### Development Mode
```powershell
.\build.ps1 dev
```
Starts Wails dev server with hot reload

### Clean Build Artifacts
```powershell
.\build.ps1 clean
```
Removes `build/` and `frontend/dist/` directories

### Run Tests
```powershell
.\build.ps1 test
```
Executes Go test suite

### Show Help
```powershell
.\build.ps1 -Help
```

## Version Update Examples

### Auto-increment build number (2.0.0.1 → 2.0.0.2)
```
Current version: 2.0.0.1
Enter new version or press Enter to increment build: [press Enter]
```

### Set new major/minor/patch version
```
Current version: 2.0.0.1
Enter new version or press Enter to increment build: 2.1.0
→ Results in: 2.1.0.1 (build resets to 1)
```

## Build Output

Successful builds create:
- `build/bin/LaunchRDP.exe` (production executable)
- Shows file size in MB
- Contains embedded version info from `version.go`

## Architecture

```
version.go (LEADING SOURCE)
    ↓
build.ps1 (syncs on build)
    ↓
wails.json (productVersion)
    ↓
Windows manifest templates
    ↓
LaunchRDP.exe (with version metadata)
```

## Important Notes

⚠️ **Always use `build.ps1` for version updates** - manual edits may cause inconsistencies  
⚠️ **Never edit version in `wails.json`** - it's automatically generated from `version.go`  
✅ **Only edit `Version` constant in `version.go`** if manually updating versions  
✅ **Build script handles all synchronization automatically**
