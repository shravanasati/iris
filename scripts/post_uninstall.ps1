$startupFolder = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"
$shortcutPath = "$startupFolder\iris.lnk"

if (Test-Path $shortcutPath -PathType Leaf) {
	Remove-Item $shortcutPath -Force
    Write-Host "Removed iris from startup applications."
}