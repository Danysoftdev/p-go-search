#!/bin/bash

set -e

echo "📁 Desplegando microservicio p-go-search..."

# Namespace
kubectl apply -f k8s/search/namespace-search.yaml

# Secret
kubectl apply -f k8s/search/secrets-search.yaml

# Deployment
kubectl apply -f k8s/search/deployment-search.yaml

# Esperar a que esté listo
echo "⏳ Esperando a que p-go-search esté listo..."
kubectl wait --namespace=p-go-search \
  --for=condition=available deployment/search-deployment \
  --timeout=90s

# Service
kubectl apply -f k8s/search/service-search.yaml

# Ingress
kubectl apply -f k8s/search/ingress.yaml

echo "✅ p-go-search desplegado correctamente."

echo -e "\n🔍 Estado actual:"
kubectl get all -n p-go-search
kubectl get ingress -n p-go-search
