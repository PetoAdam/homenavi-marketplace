# Homenavi Marketplace Helm Chart

Local Helm deployment of Homenavi Marketplace (Postgres + API + Web + optional Nginx ingress gateway).

## Install

```bash
helm upgrade --install homenavi-marketplace ./helm/homenavi-marketplace -n homenavi-marketplace --create-namespace
```

## Validate

```bash
helm lint ./helm/homenavi-marketplace
helm template homenavi-marketplace ./helm/homenavi-marketplace > /tmp/homenavi-marketplace-rendered.yaml
```

## Local access

If `nginx.enabled=true` (default):

```bash
kubectl -n homenavi-marketplace port-forward svc/homenavi-marketplace 3010:80
```

Then open `http://127.0.0.1:3010`.
