# LaunchRDP Build Script
# Build the WebView2-based LaunchRDP application

param(
    [string]$Target = "package"
)

function Get-CurrentVersion {
    if (Test-Path "version.go") {
        $versionContent = Get-Content "version.go" -Raw
        $versionMatch = [regex]::Match($versionContent, 'Version\s*=\s*"([^"]+)"')
        
        if ($versionMatch.Success) {
            $versionString = $versionMatch.Groups[1].Value
            $parts = $versionString.Split('.')
            
            if ($parts.Length -eq 3) {
                # Format: Major.Minor.Patch (add Build=0)
                $major = [int]$parts[0]
                $minor = [int]$parts[1]
                $patch = [int]$parts[2]
                $build = 0
            }
            elseif ($parts.Length -eq 4) {
                # Format: Major.Minor.Patch.Build
                $major = [int]$parts[0]
                $minor = [int]$parts[1]
                $patch = [int]$parts[2]
                $build = [int]$parts[3]
            }
            else {
                # Invalid format, default
                $major = 1; $minor = 0; $patch = 0; $build = 0
            }
            
            return @{
                Major    = $major
                Minor    = $minor
                Patch    = $patch
                Build    = $build
                String   = "$major.$minor.$patch.$build"
                GoString = $versionString
            }
        }
    }
    
    # Fallback
    return @{
        Major    = 1
        Minor    = 0
        Patch    = 0
        Build    = 0
        String   = "1.0.0.0"
        GoString = "1.0.0.0"
    }
}

function Update-Version {
    $currentVersion = Get-CurrentVersion
    
    Write-Host ""
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "        LaunchRDP Version Manager       " -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "Current version: $($currentVersion.String)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Enter new version (Major.Minor.Patch) or press Enter to keep current and increment build:"
    $userInput = Read-Host "New version"
    
    if ([string]::IsNullOrWhiteSpace($userInput)) {
        # Just increment build number
        $newBuild = $currentVersion.Build + 1
        $newVersion = @{
            Major  = $currentVersion.Major
            Minor  = $currentVersion.Minor
            Patch  = $currentVersion.Patch
            Build  = $newBuild
            String = "$($currentVersion.Major).$($currentVersion.Minor).$($currentVersion.Patch).$newBuild"
        }
        Write-Host "Incrementing build number: $($currentVersion.String) → $($newVersion.String)" -ForegroundColor Green
    }
    else {
        # Parse new version
        $parts = $userInput.Split('.')
        if ($parts.Length -ne 3) {
            Write-Host "Invalid version format! Use Major.Minor.Patch (e.g., 1.2.3)" -ForegroundColor Red
            exit 1
        }
        
        try {
            $major = [int]$parts[0]
            $minor = [int]$parts[1]
            $patch = [int]$parts[2]
            $newVersion = @{
                Major  = $major
                Minor  = $minor
                Patch  = $patch
                Build  = 1
                String = "$major.$minor.$patch.1"
            }
            Write-Host "Updating version: $($currentVersion.String) → $($newVersion.String)" -ForegroundColor Green
        }
        catch {
            Write-Host "Invalid version numbers! Use integers only." -ForegroundColor Red
            exit 1
        }
    }
    
    # First update version.go (source of truth)
    if (Test-Path "version.go") {
        $goContent = Get-Content "version.go" -Raw
        $newGoVersion = "$($newVersion.Major).$($newVersion.Minor).$($newVersion.Patch).$($newVersion.Build)"
        $updatedGoContent = $goContent -replace 'Version\s*=\s*"[^"]+"', "Version    = `"$newGoVersion`""
        $updatedGoContent | Set-Content "version.go" -NoNewline
        Write-Host "✓ Updated version.go to $newGoVersion" -ForegroundColor Green
    }
    
    # Then update versioninfo.json (for Windows resources)
    if (Test-Path "versioninfo.json") {
        $versionInfo = Get-Content "versioninfo.json" | ConvertFrom-Json
        $versionInfo.FixedFileInfo.FileVersion.Major = $newVersion.Major
        $versionInfo.FixedFileInfo.FileVersion.Minor = $newVersion.Minor
        $versionInfo.FixedFileInfo.FileVersion.Patch = $newVersion.Patch
        $versionInfo.FixedFileInfo.FileVersion.Build = $newVersion.Build
        $versionInfo.FixedFileInfo.ProductVersion.Major = $newVersion.Major
        $versionInfo.FixedFileInfo.ProductVersion.Minor = $newVersion.Minor
        $versionInfo.FixedFileInfo.ProductVersion.Patch = $newVersion.Patch
        $versionInfo.FixedFileInfo.ProductVersion.Build = $newVersion.Build
        $versionInfo.StringFileInfo.FileVersion = $newVersion.String
        $versionInfo.StringFileInfo.ProductVersion = $newVersion.String
        
        $versionInfo | ConvertTo-Json -Depth 10 | Set-Content "versioninfo.json"
        Write-Host "Updated versioninfo.json" -ForegroundColor Green
        
        # Update manifest file (only main assemblyIdentity version, not Common-Controls)
        if (Test-Path "res\LaunchRDP.manifest") {
            $manifestContent = Get-Content "res\LaunchRDP.manifest" -Raw
            # Match the main assemblyIdentity block (first one) specifically
            $pattern = '(<assemblyIdentity\s+version=")[^"]*("\s+processorArchitecture="\*"\s+name="LaunchRDP")'
            $replacement = "`${1}$($newVersion.String)`${2}"
            $updatedManifest = $manifestContent -replace $pattern, $replacement
            $updatedManifest | Set-Content "res\LaunchRDP.manifest" -NoNewline
            Write-Host "✓ Updated manifest to $($newVersion.String)" -ForegroundColor Green
        }
        
        # Regenerate resource file
        Write-Host "Regenerating Windows resources..." -ForegroundColor Cyan
        goversioninfo versioninfo.json
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Resources updated successfully" -ForegroundColor Green
        }
        else {
            Write-Host "Failed to update resources" -ForegroundColor Red
            exit 1
        }
    }
    
    return $newVersion
}

function Invoke-Package {
    $version = Update-Version
    
    Write-Host ""
    Write-Host "Building LaunchRDP v$($version.String)..." -ForegroundColor Cyan
    
    # Clean up modules
    go mod tidy
    
    # Build WebView2-only executable
    go build -ldflags "-s -w -H windowsgui" -o "LaunchRDP.exe" .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "═══════════════════════════════════════" -ForegroundColor Green
        Write-Host "          BUILD SUCCESSFUL!            " -ForegroundColor Green
        Write-Host "═══════════════════════════════════════" -ForegroundColor Green
        $size = (Get-Item "LaunchRDP.exe").Length
        $sizeMB = [math]::Round($size / 1MB, 2)
        Write-Host "✓ LaunchRDP.exe v$($version.String) created" -ForegroundColor Green
        Write-Host "✓ Executable size: $sizeMB MB" -ForegroundColor Cyan
        Write-Host "✓ WebView2-only (no external dependencies)" -ForegroundColor Cyan
        Write-Host "✓ Windows resources embedded (icon, version info)" -ForegroundColor Cyan
    }
    else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Invoke-Clean {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Cyan
    if (Test-Path "LaunchRDP.exe") {
        Remove-Item "LaunchRDP.exe"
        Write-Host "✓ Removed LaunchRDP.exe" -ForegroundColor Green
    }
    if (Test-Path "resource.syso") {
        Remove-Item "resource.syso"
        Write-Host "✓ Removed resource.syso" -ForegroundColor Green
    }
    Write-Host "Clean complete." -ForegroundColor Green
}

function Invoke-Run {
    Invoke-Package
    if (Test-Path "LaunchRDP.exe") {
        Write-Host ""
        Write-Host "Starting LaunchRDP..." -ForegroundColor Cyan
        .\LaunchRDP.exe
    }
}

function Invoke-Test {
    Write-Host "Running tests..." -ForegroundColor Cyan
    go test .\...
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ All tests passed" -ForegroundColor Green
    }
    else {
        Write-Host "✗ Tests failed" -ForegroundColor Red
    }
}

function Invoke-Format {
    Write-Host "Formatting Go code..." -ForegroundColor Cyan
    go fmt .\...
    Write-Host "✓ Code formatted" -ForegroundColor Green
}

function Show-Help {
    Write-Host ""
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "           LaunchRDP Build Tool         " -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [command]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor White
    Write-Host "  package  - Build WebView2 executable with version management (default)" -ForegroundColor Green
    Write-Host "  clean    - Remove build artifacts" -ForegroundColor Yellow
    Write-Host "  run      - Build and run the application" -ForegroundColor Cyan
    Write-Host "  test     - Run all tests" -ForegroundColor Magenta
    Write-Host "  fmt      - Format source code" -ForegroundColor Blue
    Write-Host "  help     - Show this help message" -ForegroundColor White
    Write-Host ""
    Write-Host "Features:" -ForegroundColor White
    Write-Host "  ✓ Interactive version management" -ForegroundColor Green
    Write-Host "  ✓ Automatic Windows resource embedding" -ForegroundColor Green
    Write-Host "  ✓ WebView2-only build (no external dependencies)" -ForegroundColor Green
    Write-Host "  ✓ Icon and metadata integration" -ForegroundColor Green
    Write-Host ""
}

switch ($Target.ToLower()) {
    "package" { Invoke-Package }
    "clean" { Invoke-Clean }
    "run" { Invoke-Run }
    "test" { Invoke-Test }
    "fmt" { Invoke-Format }
    "format" { Invoke-Format }
    "help" { Show-Help }
    default { Invoke-Package }  # Default to package
}