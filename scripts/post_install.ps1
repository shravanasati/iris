$executablePath = "$dir\iris.exe"
$startupFolder = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"
$shortcutPath = "$startupFolder\iris.lnk"

if (Test-Path $shortcutPath -PathType Leaf) {
    Exit 0
}

$confirmation = read-host "Would you like to add iris to startup applications? (y/N) "

if ($confirmation -eq "Y" -or $confirmation -eq "y") {

    $WshShell = New-Object -comObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut($shortcutPath)
    $Shortcut.TargetPath = $executablePath
    $Shortcut.Save()

    Write-Host "iris has been added to the list of startup applications."
}
else {
    write-host "OK."
}