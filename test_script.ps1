Start-Job -ScriptBlock {
    ./reset.ps1
} 

Start-Sleep -Seconds 60
Write-Host " Lancement des tests unitaires"
cd app
go test -v