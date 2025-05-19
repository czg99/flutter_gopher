@echo off
setlocal enabledelayedexpansion

cd /d %~dp0

where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go compiler not found. Please install Go.
    exit /b 1
)

where zig >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Zig compiler not found. Please install Zig.
    exit /b 1
)

set OUTPUT_NAME={{.LibName}}
set OUTPUT_FILE=%OUTPUT_NAME%.dll
set OUTPUT_DIR=%~dp0
set GO_SRC=../src

echo Detecting system architecture...
set GOARCH=amd64
set ZIG_TARGET=x86_64-windows-gnu

:: Detecting ARM architecture
reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PROCESSOR_ARCHITECTURE | findstr "ARM" >nul
if %ERRORLEVEL% EQU 0 (
    set GOARCH=arm64
    set ZIG_TARGET=aarch64-windows-gnu
    echo ARM architecture detected
) else (
    echo x86_64 architecture detected
)

echo Compiling Go code to shared library...

set CGO_ENABLED=1
set GOOS=windows
set GOARCH=!GOARCH!
set CC=zig cc -target !ZIG_TARGET!
set CXX=zig c++ -target !ZIG_TARGET!

go build -C %GO_SRC% -ldflags "-s -w" -trimpath -buildmode=c-shared -o "%OUTPUT_DIR%%OUTPUT_FILE%"

if %ERRORLEVEL% NEQ 0 (
    echo Error: Go compilation failed, error code: %ERRORLEVEL%
    echo Please check Go source code or compilation environment settings
    exit /b %ERRORLEVEL%
) else (
    echo Compiled successfully: %OUTPUT_DIR%%OUTPUT_FILE%
)

if exist "%OUTPUT_DIR%%OUTPUT_NAME%.h" del "%OUTPUT_DIR%%OUTPUT_NAME%.h"

echo Build process completed
exit /b 0