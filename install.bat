@echo off
setlocal enabledelayedexpansion

set BINARY_NAME=pangea
set INSTALL_PATH=C:\Program Files\

echo Running Pangea-CLI install script...
echo OS: %OS%
echo Arch: %PROCESSOR_ARCHITECTURE%

REM Check if the file named "pangea" exists in the current directory
if not exist "%BINARY_NAME%" (
    echo There is no binary called '%BINARY_NAME%' in this folder.
    pause
    exit /b 1
)

REM Copy the binary to the install path
copy "%BINARY_NAME%" "%INSTALL_PATH%" >nul 2>&1
if errorlevel 1 (
    echo Failed to copy the binary. Please run this script as an administrator.
    pause
    exit /b 1
)

echo Installation success.
pause
