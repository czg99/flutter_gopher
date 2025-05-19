Set-Location -Path $PSScriptRoot

if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Error "Error: Go compiler not found. Please install Go."
    exit 1
}

if (-not (Get-Command "zig" -ErrorAction SilentlyContinue)) {
    Write-Error "Error: Zig compiler not found. Please install Zig."
    exit 1
}

$OUTPUT_NAME = "{{.LibName}}"
$OUTPUT_FILE = "$OUTPUT_NAME.dll"
$OUTPUT_DIR = $PSScriptRoot
$GO_SRC = "../src"
$TIMESTAMP_FILE = ".last_build_time"

# Function to check if source code has been updated
function Check-SourceChanges {
    if (-not (Test-Path -Path $TIMESTAMP_FILE)) {
        return $true
    }

    if (-not (Test-Path -Path (Join-Path -Path $OUTPUT_DIR -ChildPath $OUTPUT_FILE))) {
        return $true
    }
    
    $lastBuildTime = [datetime]::Parse((Get-Content -Path $TIMESTAMP_FILE))

    $newestFileTime = Get-ChildItem -Path $GO_SRC -Filter "*.go" -Recurse | 
                      Select-Object -ExpandProperty LastWriteTime | 
                      Sort-Object -Descending | 
                      Select-Object -First 1
    
    if ($null -eq $newestFileTime) {
        return $true
    }
    
    return $newestFileTime -gt $lastBuildTime
}

# Function to save current build timestamp
function Save-BuildTime {
    Get-Date | Out-File -FilePath $TIMESTAMP_FILE
}

# Check if source code has been updated
$needCompile = Check-SourceChanges
if (-not $needCompile) {
    Write-Host "Source code unchanged, skipping compilation"
    exit 0
}

Write-Host "Detecting system architecture..."
$GOARCH = "amd64"
$ZIG_TARGET = "x86_64-windows-gnu"

$processorArchitecture = (Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" -Name PROCESSOR_ARCHITECTURE).PROCESSOR_ARCHITECTURE
if ($processorArchitecture -match "ARM") {
    $GOARCH = "arm64"
    $ZIG_TARGET = "aarch64-windows-gnu"
    Write-Host "ARM architecture detected"
} else {
    Write-Host "x86_64 architecture detected"
}

Write-Host "Compiling Go code to shared library..."

$env:CGO_ENABLED = 1
$env:GOOS = "windows"
$env:GOARCH = $GOARCH
$env:CC = "zig cc -target $ZIG_TARGET"
$env:CXX = "zig c++ -target $ZIG_TARGET"

& go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-shared -o "$OUTPUT_DIR\$OUTPUT_FILE"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Go compilation failed, error code: $LASTEXITCODE"
    Write-Error "Please check Go source code or compilation environment settings"
    exit 1
}

Write-Host "Compiled successfully: $OUTPUT_DIR\$OUTPUT_FILE"

$headerFile = Join-Path -Path $OUTPUT_DIR -ChildPath "$OUTPUT_NAME.h"
if (Test-Path -Path $headerFile) {
    Remove-Item -Path $headerFile -Force
}

# Save current build timestamp
Save-BuildTime

Write-Host "Build process completed"