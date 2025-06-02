@echo off

cd /d %~dp0

powershell.exe -ExecutionPolicy Bypass -File build.ps1