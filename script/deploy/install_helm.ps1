# Installer helm
scoop install helm

# Créer un chart Helm
helm create pocnokc

# Ajouter les dépôts Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Installer Prometheus dans le namespace monitoring (le crée s'il n'existe pas)
helm install prometheus prometheus-community/prometheus `
  --namespace monitoring --create-namespace

# Installer Grafana dans le même namespace
helm install grafana grafana/grafana `
  --namespace monitoring

# Récupérer le mot de passe admin de Grafana
$secret = kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}"
$decodedPassword = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($secret))
Write-Output "Grafana admin password: $decodedPassword"
