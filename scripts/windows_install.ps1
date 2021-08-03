Write-Host "Downloading iris..."

$url = "https://github.com/Shravan-1908/iris/releases/latest/download/iris-windows-amd64.exe"

$dir = $env:USERPROFILE + "\.iris"
$filepath = $env:USERPROFILE + "\.iris\iris.exe"

try {
    [System.IO.Directory]::CreateDirectory($dir)
}
catch {
    Write-Host "Failed to create directory!"
    [Environment]::Exit(1)
}


try {
    (Invoke-WebRequest -Uri $url -OutFile $filepath)
}
catch {
    Write-Host "Failed to download the executable file! Make sure youve active internet connection."
    [Environment]::Exit(1)
}

try {
    Write-Host "Adding iris to PATH..."
    [Environment]::SetEnvironmentVariable(
        "Path",
        [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine) + ";"+$dir,
        [EnvironmentVariableTarget]::Machine)
}
catch {
    Write-Host "Failed to add iris to PATH! Make sure you've opened powershell as Admin."
    [Environment]::Exit(1)
}

try {
    Write-Host "Adding iris to startup applications..."

    $startuppath = $env:USERPROFILE + "\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\iris.lnk"

    $WshShell = New-Object -comObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut($startuppath)
    $Shortcut.TargetPath = $filepath
    $Shortcut.Save()
}
catch {
    Write-Host "Failed to add iris to startup applications! Make sure you've opened powershell as Admin."
    [Environment]::Exit(1)
}

Write-Host "iris installation is successful!"