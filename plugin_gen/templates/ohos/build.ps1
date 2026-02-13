# Receive parameters: OHOS NDK home path and minimum API version
param(
    [Parameter(Mandatory=$true)]
    [string]$OHOS_NDK_HOME
)

Write-Host "OHOS_NDK_HOME: $OHOS_NDK_HOME"

Set-Location -Path $PSScriptRoot

if ([string]::IsNullOrEmpty($OHOS_NDK_HOME)) {
    Write-Error "Error: Please provide OHOS NDK home path as the first parameter"
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
$GO_SRC = "../gosrc"
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
    
    $archs = @("arm64-v8a", "armeabi-v7a", "x86_64")
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

$CC = Join-Path -Path $OHOS_NDK_HOME -ChildPath "native\llvm\bin\clang.exe"
$CXX = Join-Path -Path $OHOS_NDK_HOME -ChildPath "native\llvm\bin\clang++.exe"

$archs = @("arm64-v8a", "armeabi-v7a", "x86_64")
foreach ($arch in $archs) {
    Write-Host "Compiling $arch architecture..."
    
    switch ($arch) {
        "arm64-v8a" {
            $env:GOARCH = "arm64"
            $CC_TARGET = "aarch64-linux-ohos"
        }
        "armeabi-v7a" {
            $env:GOARCH = "arm"
            $env:GOARM = "7"
            $CC_TARGET = "armv7-linux-ohos"
        }
        "x86_64" {
            $env:GOARCH = "amd64"
            $CC_TARGET = "x86_64-linux-ohos"
        }
    }
    
    $archDir = Join-Path -Path $OUTPUT_DIR -ChildPath $arch
    if (-not (Test-Path -Path $archDir)) {
        New-Item -Path $archDir -ItemType Directory | Out-Null
    }
    
    $env:CC = "$CC --target=$CC_TARGET --sysroot=$OHOS_NDK_HOME/native/sysroot"
    $env:CXX = "$CXX --target=$CC_TARGET --sysroot=$OHOS_NDK_HOME/native/sysroot"

    $env:CGO_CFLAGS = "-I${PWD}/log-adaptor/include"
    $env:CGO_LDFLAGS = "-L${PWD}/log-adaptor/dist/${arch}"
    
    $outputPath = Join-Path -Path $archDir -ChildPath $OUTPUT_FILE

    & go build -C $GO_SRC -ldflags "-s -w -extldflags '-Wl,-soname,$OUTPUT_FILE'" -trimpath -buildmode=c-shared -o "$outputPath"
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