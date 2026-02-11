package db

import (
	"context"

	"gorm.io/gorm"
)

func Migrate(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).AutoMigrate(
		&Integration{},
		&IntegrationDownloadEvent{},
	); err != nil {
		return err
	}

	// GORM does not support partial unique indexes; use raw SQL for latest-only uniqueness.
	if err := db.WithContext(ctx).Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS integrations_listen_path_latest_unique
  ON integrations (listen_path)
  WHERE latest = TRUE
`).Error; err != nil {
		return err
	}

	// Enforce name uniqueness for latest integrations (partial unique index).
	if err := db.WithContext(ctx).Exec(`
CREATE UNIQUE INDEX IF NOT EXISTS integrations_name_latest_unique
  ON integrations (name)
  WHERE latest = TRUE
`).Error; err != nil {
		return err
	}

	return nil
}
