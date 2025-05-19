@echo off
setlocal enabledelayedexpansion

cd /d %~dp0

:: Receive parameters: NDK path and minimum API version
set NDK_PATH=%1
set MIN_API=%2

if "%NDK_PATH%"=="" (
    echo Error: Please provide NDK path as the first parameter
    exit /b 1
)

if "%MIN_API%"=="" (
    echo Error: Please provide minimum API version as the second parameter
    exit /b 1
)

where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Go compiler not found. Please install Go.
    exit /b 1
)

set OUTPUT_NAME={{.LibName}}
set OUTPUT_FILE=lib%OUTPUT_NAME%.so
set OUTPUT_HEADER=lib%OUTPUT_NAME%.h
set OUTPUT_DIR=%CD%\libs
set GO_SRC=..\src

if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

set CGO_ENABLED=1
set GOOS=android

set HOST_OS=windows
set CC="%NDK_PATH%\toolchains\llvm\prebuilt\%HOST_OS%-x86_64\bin\clang.exe"
set CXX="%NDK_PATH%\toolchains\llvm\prebuilt\%HOST_OS%-x86_64\bin\clang++.exe"

set ARCHS=arm64-v8a armeabi-v7a x86 x86_64
for %%A in (%ARCHS%) do (
    echo Compiling %%A architecture...
    
    if "%%A"=="arm64-v8a" (
        set GOARCH=arm64
        set CC_TARGET=aarch64-linux-android%MIN_API%
    ) else if "%%A"=="armeabi-v7a" (
        set GOARCH=arm
        set GOARM=7
        set CC_TARGET=armv7a-linux-androideabi%MIN_API%
    ) else if "%%A"=="x86" (
        set GOARCH=386
        set CC_TARGET=i686-linux-android%MIN_API%
    ) else if "%%A"=="x86_64" (
        set GOARCH=amd64
        set CC_TARGET=x86_64-linux-android%MIN_API%
    )
    
    if not exist "%OUTPUT_DIR%\%%A" mkdir "%OUTPUT_DIR%\%%A"
    
    set CC=%CC% --target=!CC_TARGET!
    set CXX=%CXX% --target=!CC_TARGET!
    
    go build -C %GO_SRC% -ldflags "-s -w" -trimpath -buildmode=c-shared -o "%OUTPUT_DIR%\%%A\%OUTPUT_FILE%"
    
    if %ERRORLEVEL% neq 0 (
        echo Error: Go compilation failed for architecture %%A
        exit /b 1
    )
    
    if exist "%OUTPUT_DIR%\%%A\%OUTPUT_HEADER%" del "%OUTPUT_DIR%\%%A\%OUTPUT_HEADER%"
    
    echo %%A compiled successfully: %OUTPUT_DIR%\%%A\%OUTPUT_FILE%
)

echo Build process completed
exit /b 0
