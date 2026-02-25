@echo off
setlocal

set "APP=%~dp0download-cleaner.exe"

if not exist "%APP%" (
    echo Error: download-cleaner.exe was not found in this folder.
    echo Keep this .bat file and the .exe in the same folder.
    echo.
    pause
    exit /b 1
)

echo Download Cleaner
echo ----------------
echo.
echo Press Y when asked to proceed, or N to cancel.
echo.

"%APP%"
set "EXITCODE=%ERRORLEVEL%"

echo.
if "%EXITCODE%"=="0" (
    echo Finished successfully.
) else (
    echo Finished with errors. Exit code: %EXITCODE%
)
echo.
echo Press any key to close this window.
pause >nul
exit /b %EXITCODE%
