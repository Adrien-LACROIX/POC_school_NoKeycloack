#!/bin/bash

set -e

# Vérifie si le script est exécuté avec les privilèges root
if [ "$EUID" -ne 0 ]; then
    echo "Veuillez exécuter ce script en tant que root (sudo)."
    exit 1
fi

echo "Mise à jour des paquets..."
apt update -y
apt upgrade -y

echo "Suppression des anciennes versions de Docker (si présentes)..."
apt remove -y docker docker-engine docker.io containerd runc || true

echo "Installation des dépendances..."
apt install -y ca-certificates curl gnupg lsb-release

echo "Ajout de la clé GPG officielle de Docker..."
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
    gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo "Ajout du dépôt Docker au sources.list..."
ARCH=$(dpkg --print-architecture)
RELEASE=$(lsb_release -cs)
echo \
  "deb [arch=$ARCH signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $RELEASE stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "Mise à jour des dépôts avec Docker inclus..."
apt update -y

echo "Installation de Docker Engine, CLI, containerd et Docker Compose plugin..."
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "Activation et démarrage de Docker..."
systemctl enable docker
systemctl start docker

echo "Ajout de l'utilisateur actuel au groupe docker (nécessite une reconnexion)..."
usermod -aG docker "$USER"


echo "🛑 Arrêt et suppression des conteneurs + volumes..."
docker-compose down -v

echo "🚀 Reconstruction et démarrage des conteneurs..."
docker-compose up --build -d

echo "▶️ Lancement de l'application..."
go run ./main.go &

echo "⏳ Attente du démarrage de l'application (15s)..."
sleep 15

echo "🧪 Lancement des tests unitaires..."
go test -v

