Write-Host "Arrêt et suppression des conteneurs + volumes..."
docker-compose down -v

Write-Host "Reconstruction et démarrage des conteneurs..."
docker-compose up --build -d

Write-Host "Lancement de l application..."
Start-Process powershell -ArgumentList "go run ./main.go"
Start-Sleep -Seconds 15

Write-Host " Lancement des tests unitaires"
go test -v