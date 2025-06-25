# Receive parameters: NDK path and minimum API version
param(
    [Parameter(Mandatory=$true)]
    [string]$NDK_PATH,
    
    [Parameter(Mandatory=$true)]
    [string]$MIN_API
)

Write-Host "NDK_PATH: $NDK_PATH"
Write-Host "MIN_API: $MIN_API"

Set-Location -Path $PSScriptRoot

if ([string]::IsNullOrEmpty($NDK_PATH)) {
    Write-Error "Error: Please provide NDK path as the first parameter"
    exit 1
}

if ([string]::IsNullOrEmpty($MIN_API)) {
    Write-Error "Error: Please provide minimum API version as the second parameter"
    exit 1
}

# Check if Go is installed
if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Error "Error: Go compiler not found. Please install Go."
    exit 1
}

$OUTPUT_NAME = "{{.LibName}}"
$OUTPUT_FILE = "lib$OUTPUT_NAME.so"
$OUTPUT_HEADER = "lib$OUTPUT_NAME.h"
$OUTPUT_DIR = Join-Path -Path $PSScriptRoot -ChildPath "libs"
$GO_SRC = "../src"
$TIMESTAMP_FILE = ".last_build_time"

if (-not (Test-Path -Path $OUTPUT_DIR)) {
    New-Item -Path $OUTPUT_DIR -ItemType Directory | Out-Null
}

# Function to check if source code has been updated
function Check-SourceChanges {
    if (-not (Test-Path -Path $TIMESTAMP_FILE)) {
        return $true
    }
    
    $lastBuildTime = [datetime]::Parse((Get-Content -Path $TIMESTAMP_FILE))
    
    $archs = @("arm64-v8a", "armeabi-v7a", "x86", "x86_64")
    foreach ($arch in $archs) {
        $archOutputFile = Join-Path -Path $OUTPUT_DIR -ChildPath "$arch\$OUTPUT_FILE"
        if (-not (Test-Path -Path $archOutputFile)) {
            return $true
        }
    }
    
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

$env:CGO_ENABLED = 1
$env:GOOS = "android"

$HOST_OS = "windows"
$CC = Join-Path -Path $NDK_PATH -ChildPath "toolchains\llvm\prebuilt\$HOST_OS-x86_64\bin\clang.exe"
$CXX = Join-Path -Path $NDK_PATH -ChildPath "toolchains\llvm\prebuilt\$HOST_OS-x86_64\bin\clang++.exe"

$archs = @("arm64-v8a", "armeabi-v7a", "x86", "x86_64")
foreach ($arch in $archs) {
    Write-Host "Compiling $arch architecture..."
    
    switch ($arch) {
        "arm64-v8a" {
            $env:GOARCH = "arm64"
            $CC_TARGET = "aarch64-linux-android$MIN_API"
        }
        "armeabi-v7a" {
            $env:GOARCH = "arm"
            $env:GOARM = "7"
            $CC_TARGET = "armv7a-linux-androideabi$MIN_API"
        }
        "x86" {
            $env:GOARCH = "386"
            $CC_TARGET = "i686-linux-android$MIN_API"
        }
        "x86_64" {
            $env:GOARCH = "amd64"
            $CC_TARGET = "x86_64-linux-android$MIN_API"
        }
    }
    
    $archDir = Join-Path -Path $OUTPUT_DIR -ChildPath $arch
    if (-not (Test-Path -Path $archDir)) {
        New-Item -Path $archDir -ItemType Directory | Out-Null
    }
    
    $env:CC = "$CC --target=$CC_TARGET"
    $env:CXX = "$CXX --target=$CC_TARGET"
    
    $outputPath = Join-Path -Path $archDir -ChildPath $OUTPUT_FILE

    & go build -C $GO_SRC -ldflags "-s -w" -trimpath -buildmode=c-shared -o "$outputPath"
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Error: Go compilation failed for architecture $arch"
        exit 1
    }
    
    $headerPath = Join-Path -Path $archDir -ChildPath $OUTPUT_HEADER
    if (Test-Path -Path $headerPath) {
        Remove-Item -Path $headerPath -Force
    }
    
    Write-Host "$arch compiled successfully: $outputPath"
}

# Save current build timestamp
Save-BuildTime

Write-Host "Build process completed"