Write-Host "Lancement de Docker Desktop..."
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
Start-Sleep -Seconds 15

cd ../deployments
Write-Host "Arrêt et suppression des conteneurs + volumes..."
docker-compose down -v

Write-Host "Reconstruction et démarrage des conteneurs..."
docker-compose up --build -d
