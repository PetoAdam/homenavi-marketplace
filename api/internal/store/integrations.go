package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrListenPathInUse = errors.New("listen_path already in use")
var ErrNameInUse = errors.New("name already in use")

func ListIntegrations(ctx context.Context, pool *pgxpool.Pool, latestOnly bool) ([]models.Integration, error) {
	query := `
SELECT id, name, version, description, manifest_url, manifest, image, images, assets, listen_path,
       repo_url, release_tag, publisher, verified, latest, created_at, updated_at
FROM integrations`
	if latestOnly {
		query += " WHERE latest = true"
	}
	query += " ORDER BY name ASC"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Integration{}
	for rows.Next() {
		item, err := scanIntegration(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func GetIntegration(ctx context.Context, pool *pgxpool.Pool, id string, version string) (*models.Integration, error) {
	query := `
SELECT id, name, version, description, manifest_url, manifest, image, images, assets, listen_path,
       repo_url, release_tag, publisher, verified, latest, created_at, updated_at
FROM integrations
WHERE id = $1`
	args := []any{id}
	if version != "" {
		query += " AND version = $2"
		args = append(args, version)
	} else {
		query += " AND latest = true"
	}
	row := pool.QueryRow(ctx, query, args...)
	item, err := scanIntegration(row)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func ListVersions(ctx context.Context, pool *pgxpool.Pool, id string) ([]models.Integration, error) {
	query := `
SELECT id, name, version, description, manifest_url, manifest, image, images, assets, listen_path,
       repo_url, release_tag, publisher, verified, latest, created_at, updated_at
FROM integrations
WHERE id = $1
ORDER BY created_at DESC`

	rows, err := pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Integration{}
	for rows.Next() {
		item, err := scanIntegration(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func PublishIntegration(ctx context.Context, pool *pgxpool.Pool, req models.PublishRequest, verified bool) (*models.Integration, error) {
	if req.ListenPath == "" {
		return nil, errors.New("listen_path is required")
	}

	log.Printf("store publish integration id=%q version=%q listen_path=%q verified=%t", req.ID, req.Version, req.ListenPath, verified)

	if err := ensureListenPathAvailable(ctx, pool, req.ListenPath, req.ID); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(ctx, pool, req.Name, req.ID); err != nil {
		return nil, err
	}

	manifestJSON, err := json.Marshal(req.Manifest)
	if err != nil {
		return nil, fmt.Errorf("manifest json invalid: %w", err)
	}
	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		return nil, fmt.Errorf("images json invalid: %w", err)
	}
	assetsJSON, err := json.Marshal(req.Assets)
	if err != nil {
		return nil, fmt.Errorf("assets json invalid: %w", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, "UPDATE integrations SET latest = false WHERE id = $1", req.ID); err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO integrations (
  id, version, name, description, manifest_url, manifest, image, images, assets, listen_path,
  repo_url, release_tag, publisher, verified, latest
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,true
)
ON CONFLICT (id, version) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  manifest_url = EXCLUDED.manifest_url,
  manifest = EXCLUDED.manifest,
  image = EXCLUDED.image,
  images = EXCLUDED.images,
  assets = EXCLUDED.assets,
  listen_path = EXCLUDED.listen_path,
  repo_url = EXCLUDED.repo_url,
  release_tag = EXCLUDED.release_tag,
  publisher = EXCLUDED.publisher,
  verified = EXCLUDED.verified,
  latest = true,
  updated_at = NOW()
`,
		req.ID,
		req.Version,
		req.Name,
		req.Description,
		req.ManifestURL,
		manifestJSON,
		req.Image,
		imagesJSON,
		assetsJSON,
		req.ListenPath,
		req.RepoURL,
		req.ReleaseTag,
		req.Publisher,
		verified,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	item, err := GetIntegration(ctx, pool, req.ID, req.Version)
	if err != nil {
		return nil, err
	}
	log.Printf("store publish persisted id=%q version=%q latest=%t verified=%t", item.ID, item.Version, item.Latest, item.Verified)
	return item, nil
}

func ensureListenPathAvailable(ctx context.Context, pool *pgxpool.Pool, listenPath, id string) error {
	var existingID string
	row := pool.QueryRow(ctx, "SELECT id FROM integrations WHERE listen_path = $1 AND latest = true AND id <> $2", listenPath, id)
	scanErr := row.Scan(&existingID)
	if scanErr == nil {
		return ErrListenPathInUse
	}
	return nil
}

func ensureNameAvailable(ctx context.Context, pool *pgxpool.Pool, name, id string) error {
	var existingID string
	row := pool.QueryRow(ctx, "SELECT id FROM integrations WHERE name = $1 AND latest = true AND id <> $2", name, id)
	scanErr := row.Scan(&existingID)
	if scanErr == nil {
		return ErrNameInUse
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanIntegration(row rowScanner) (models.Integration, error) {
	var item models.Integration
	var manifestJSON []byte
	var imagesJSON []byte
	var assetsJSON []byte

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Version,
		&item.Description,
		&item.ManifestURL,
		&manifestJSON,
		&item.Image,
		&imagesJSON,
		&assetsJSON,
		&item.ListenPath,
		&item.RepoURL,
		&item.ReleaseTag,
		&item.Publisher,
		&item.Verified,
		&item.Latest,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, err
	}
	if len(manifestJSON) > 0 {
		_ = json.Unmarshal(manifestJSON, &item.Manifest)
	}
	if len(imagesJSON) > 0 {
		_ = json.Unmarshal(imagesJSON, &item.Images)
	}
	if len(assetsJSON) > 0 {
		_ = json.Unmarshal(assetsJSON, &item.Assets)
	}
	return item, nil
}
