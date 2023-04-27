$startupFolder = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"
$scriptPath = "$startupFolder\iris.bat"

if (Test-Path $scriptPath -PathType Leaf) {
	Remove-Item $scriptPath -Force
    Write-Host "Removed iris from startup applications."
}