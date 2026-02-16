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
  "compose_file": "https://.../compose/docker-compose.integration.yml",
  "repo_url": "https://github.com/PetoAdam/homenavi-spotify",
  "release_tag": "v0.1.3",
  "publisher": "Homenavi"
}
```

Validation:

- `id`, `name`, `version`, `listen_path`, `manifest_url`, `image` are required.
- `compose_file` is required (must be a URL to docker-compose.integration.yml, not dev).
- `images` max 5.
- `listen_path` must be unique across latest releases.
- `version` and `release_tag` must match the Git tag.
- `repo_url` must match the GitHub repository from the OIDC token.
- `manifest_url` must reference the same repository + tag.

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
