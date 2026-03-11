# Homenavi Marketplace

A fullstack marketplace for Homenavi integrations.

- Go REST API + Postgres
- Next.js + Tailwind frontend
- Secure publish endpoint for CI pipelines
- Validation for unique listen paths

## Project layout

```
api/   # Go REST API
web/   # Next.js marketplace UI
```

## Requirements

- Go 1.22+
- Node 20+
- Docker (optional)

## Local dev (Docker Compose)

```bash
docker compose up -d
```

- Web + API (via nginx): http://localhost:3010

## Local dev (manual)

```bash
# start db
createdb homenavi_marketplace

# api
cd api
export $(cat .env | xargs)
go run ./cmd/server

# web (new terminal)
cd web
cp .env.example .env
npm install
npm run dev
```

## API

### Health

`GET /api/health`

### List integrations

`GET /api/integrations?latest=true`

### Get integration

`GET /api/integrations/{id}`

Optional `version` query param:

`GET /api/integrations/{id}?version=v0.1.0`

### List versions

`GET /api/integrations/{id}/versions`

### Publish integration (CI only, OIDC)

`POST /api/integrations/publish-oidc`

Headers:

- `Authorization: Bearer <github-oidc-token>`

Body:

```json
{
  "id": "spotify",
  "name": "Spotify",
  "version": "v0.1.3",
  "description": "Play, pause, and browse Spotify.",
  "manifest_url": "https://raw.githubusercontent.com/PetoAdam/homenavi-spotify/v0.1.3/manifest/homenavi-integration.json",
  "manifest": {},
  "image": "ghcr.io/petoadam/homenavi-spotify:latest",
  "images": ["https://.../hero.png"],
  "assets": {"icon": "https://.../icon.svg"},
  "listen_path": "/integrations/spotify",
  "deployment_artifacts": {
    "compose": {
      "file": "https://.../compose/docker-compose.integration.yml"
    },
    "helm": {
      "chart_ref": "oci://ghcr.io/petoadam/homenavi-spotify",
      "version": "v0.1.3"
    }
  },
  "repo_url": "https://github.com/PetoAdam/homenavi-spotify",
  "release_tag": "v0.1.3",
  "publisher": "Homenavi"
}
```

Validation:

- `id`, `name`, `version`, `listen_path`, `manifest_url`, `image` are required.
- At least one deployment artifact is required:
  - `deployment_artifacts.compose.file`, or
  - `deployment_artifacts.helm.chart_ref`, or
  - `deployment_artifacts.k8s_generated.chart_ref`
- If `deployment_artifacts.compose.file` is provided, it must point to `docker-compose.integration.yml`.
- `images` max 5.
- `listen_path` must be unique across latest releases.
- `version` and `release_tag` must match the Git tag.
- `repo_url` must match the GitHub repository from the OIDC token.
- `manifest_url` must reference the same repository + tag.

## Local Minikube Helm MVP

Current MVP target is to run marketplace locally on Minikube via Helm, alongside the core Homenavi chart.

Helm chart path in this repository:

- `helm/homenavi-marketplace`

Included chart scope:

- `api` deployment + service
- `web` deployment + service
- `db` deployment + service + PVC
- `nginx` gateway service (default on)

Target install command (once chart is present):

```bash
helm upgrade --install homenavi-marketplace ./helm/homenavi-marketplace -n homenavi-marketplace --create-namespace
```

Recommended local checks:

```bash
kubectl -n homenavi-marketplace get pods
kubectl -n homenavi-marketplace get svc
kubectl -n homenavi-marketplace port-forward svc/homenavi-marketplace-api 8098:8098
curl -fsS http://127.0.0.1:8098/api/health

kubectl -n homenavi-marketplace port-forward svc/homenavi-marketplace 3010:80
# open http://127.0.0.1:3010
```

## Tests

API tests use Testcontainers and require Docker to be running.

```bash
cd api
go test ./...
```

## Release flow (integration repos)

1) Tag a release in the integration repo.
2) CI runs tests and builds the image.
3) CI requests a GitHub OIDC token and calls the marketplace publish endpoint with the release metadata.

Recommended integration CI hardening:

- Make `verify.yml` the primary quality gate (tests, `go vet`, `gosec`, Docker build, Trivy scan).
- In `release.yml`, run `verify.yml` as a required stage before publish.
- Keep central enforcement in `PetoAdam/homenavi/.github/actions/integration-release@main` (verify + `go vet` + `gosec`) so release checks cannot be bypassed by per-repo workflow edits.
- Publish signed images with SBOM + provenance.

The marketplace verifies the OIDC token and checks that the repository has a successful `verify.yml` workflow run for the tagged commit.

## Security notes

- The publish endpoint only accepts GitHub OIDC tokens.
- `listen_path` uniqueness is enforced by the API + DB index.
- Additional validation can be added in integration-proxy at runtime.
