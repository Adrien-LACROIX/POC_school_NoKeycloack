#!/bin/bash

echo "Lancement de Docker Desktop..."
cd ./deployments
# Docker Desktop sur Ubuntu se lance généralement via cette commande
systemctl --user start docker-desktop
sleep 15

echo "Arrêt et suppression des conteneurs + volumes..."
docker compose down -v

echo "Reconstruction et démarrage des conteneurs..."
docker compose up --build -d
cd ..