param(
    [string]$Target = "build",
    [switch]$Help
)

function Get-CurrentVersion {
    if (Test-Path "version.go") {
        $versionContent = Get-Content "version.go" -Raw
        $versionMatch = [regex]::Match($versionContent, 'Version\s*=\s*"([^"]+)"')
        
        if ($versionMatch.Success) {
            $versionString = $versionMatch.Groups[1].Value
            $parts = $versionString.Split('.')
            
            if ($parts.Length -eq 4) {
                $major = [int]$parts[0]; $minor = [int]$parts[1]; $patch = [int]$parts[2]; $build = [int]$parts[3]
            }
            else {
                $major = 1; $minor = 0; $patch = 0; $build = 0
            }
            
            return @{
                Major = $major; Minor = $minor; Patch = $patch; Build = $build
                String = "$major.$minor.$patch.$build"; GoString = $versionString
            }
        }
    }
    return @{ Major = 1; Minor = 0; Patch = 0; Build = 0; String = "1.0.0.0"; GoString = "1.0.0" }
}

function Update-Version {
    $currentVersion = Get-CurrentVersion
    Write-Host "Current version: $($currentVersion.GoString)" -ForegroundColor Yellow
    $userInput = Read-Host "Enter new version or press Enter to increment build"
    
    if ([string]::IsNullOrWhiteSpace($userInput)) {
        $newBuild = $currentVersion.Build + 1
        $newVersion = @{
            Major = $currentVersion.Major; Minor = $currentVersion.Minor; Patch = $currentVersion.Patch; Build = $newBuild
            String = "$($currentVersion.Major).$($currentVersion.Minor).$($currentVersion.Patch).$newBuild"
        }
    }
    else {
        $parts = $userInput.Split('.')
        if ($parts.Length -ne 3) { Write-Host "Invalid format!" -ForegroundColor Red; exit 1 }
        $major = [int]$parts[0]; $minor = [int]$parts[1]; $patch = [int]$parts[2]
        $newVersion = @{ Major = $major; Minor = $minor; Patch = $patch; Build = 1; String = "$major.$minor.$patch.1" }
    }
    
    # Update version.go (leading source)
    $goContent = Get-Content "version.go" -Raw
    $newGoVersion = "$($newVersion.Major).$($newVersion.Minor).$($newVersion.Patch).$($newVersion.Build)"
    $updatedGoContent = $goContent -replace 'Version\s*=\s*"[^"]+"', "Version    = `"$newGoVersion`""
    $updatedGoContent | Set-Content "version.go" -NoNewline
    Write-Host "Updated version.go to $newGoVersion" -ForegroundColor Green
    
    # Sync wails.json with version.go (for Windows manifest)
    $wailsJsonPath = "wails.json"
    if (Test-Path $wailsJsonPath) {
        $wailsContent = Get-Content $wailsJsonPath -Raw | ConvertFrom-Json
        $wailsVersion = "$($newVersion.Major).$($newVersion.Minor).$($newVersion.Patch)"
        $wailsContent.info.productVersion = $wailsVersion
        $wailsContent | ConvertTo-Json -Depth 10 | Set-Content $wailsJsonPath
        Write-Host "Synced wails.json productVersion to $wailsVersion" -ForegroundColor Green
    }
    
    return $newVersion
}

function Invoke-WailsBuild {
    $version = Update-Version
    Write-Host "Building LaunchRDP v$($version.String) with Wails..." -ForegroundColor Cyan
    
    go mod tidy
    wails generate
    if ($LASTEXITCODE -ne 0) { Write-Host "Failed to generate bindings" -ForegroundColor Red; exit 1 }
    
    wails build -clean
    if ($LASTEXITCODE -eq 0) {
        $exePath = "build\bin\LaunchRDP.exe"
        if (Test-Path $exePath) {
            $sizeMB = [math]::Round((Get-Item $exePath).Length / 1MB, 2)
            Write-Host "BUILD SUCCESSFUL!" -ForegroundColor Green
            Write-Host "LaunchRDP.exe v$($version.String) created ($sizeMB MB)" -ForegroundColor Green          
        }
    }
    else {
        Write-Host "Build failed!" -ForegroundColor Red; exit 1
    }
}

function Invoke-Dev { wails dev }
function Invoke-Clean { 
    if (Test-Path "build") { Remove-Item -Recurse -Force "build" }
    if (Test-Path "frontend/dist") { Remove-Item -Recurse -Force "frontend/dist" }
}
function Invoke-Test { go test ./... }

if ($Help) {
    Write-Host "Available commands: build, dev, clean, test" -ForegroundColor Cyan
    exit 0
}

switch ($Target.ToLower()) {
    { $_ -in "", "build", "package" } { Invoke-WailsBuild }
    { $_ -in "dev", "develop" } { Invoke-Dev }
    "clean" { Invoke-Clean }
    "test" { Invoke-Test }
    default { Write-Host "Unknown target: $Target" -ForegroundColor Red; exit 1 }
}
