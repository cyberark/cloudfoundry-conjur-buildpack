@echo off

pushd %USERPROFILE%\app
  for /f "tokens=1,2 delims=: "  %%a in ('.conjur\conjur-win-env.exe') do (
    for /f %%i in ('powershell -executionpolicy Unrestricted -Command "[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String('%%b'))"') do (
        set %%a=%%i
    )
  )
popd