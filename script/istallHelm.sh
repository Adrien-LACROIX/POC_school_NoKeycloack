


helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/prometheus \
  --namespace monitoring --create-namespace


helm install grafana grafana/grafana \
  --namespace monitoring


kubectl get secret --namespace monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode
