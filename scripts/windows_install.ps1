Write-Host "Downloading iris..."

$url = "https://github.com/Shravan-1908/iris/releases/latest/download/iris-windows-amd64.exe"

$dir = $env:USERPROFILE + "\.iris"
$filepath = $env:USERPROFILE + "\.iris\iris.exe"

[System.IO.Directory]::CreateDirectory($dir)
(Invoke-WebRequest -Uri $url -OutFile $filepath)

Write-Host "Adding iris to PATH..."
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine) + ";"+$dir,
    [EnvironmentVariableTarget]::Machine)

Write-Host "Adding iris to startup applications..."

$startupdir = $env:USERPROFILE + "\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup"

Write-Host "iris" >> $startupdir + "\iris.bat"

Write-Host "iris installation is successfull!"