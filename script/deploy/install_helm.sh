#!/bin/bash

echo "Installation de Helm..."
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

echo "Installation de Minikube..."
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

echo "Start Minikube and create config..."
minikube start

echo "Creation de l'application pocnokc..."
helm create pocnokc

echo "Add prometheus and grafana to repo..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

echo "test"
kubectl version --short
kubectl config current-context

echo "Install prometheus..."
helm install prometheus prometheus-community/prometheus \
  --namespace monitoring --create-namespace

echo "Install grafana..."
helm install grafana grafana/grafana \
  --namespace monitoring

echo "Dernière étape..."
kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode
