$startupFolder = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"
$scriptPath = "$startupFolder\iris.bat"

if (Test-Path $scriptPath -PathType Leaf) {
    Exit 0
}

$confirmation = read-host "Would you like to add iris to startup applications? (y/N) "

if ($confirmation -eq "Y" -or $confirmation -eq "y") {
    write-output 'powershell -NonInteractive "start-process iris -WindowStyle Hidden"' > $scriptPath
    Write-Host "iris has been added to the list of startup applications."
}
else {
    write-host "OK."
}