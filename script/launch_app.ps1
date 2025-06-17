Write-Host "Lancement de Docker Desktop..."
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
Start-Sleep -Seconds 15

cd ./deployments
Write-Host "Arrêt et suppression des conteneurs + volumes..."
docker-compose down -v

Write-Host "Reconstruction et démarrage des conteneurs..."
docker-compose up --build -d

Write-Host "Lancement de l application..."
cd ../cmd
Start-Process powershell -ArgumentList "go run ./main.go"
Start-Sleep -Seconds 15

Write-Host "Ouverture de l application dans le navigateur..."
Start-Process "http://localhost:8080"

Write-Host "Environnement pret à l emploi avec interface ouverte..."