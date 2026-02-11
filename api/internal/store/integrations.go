package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	dbmodels "github.com/PetoAdam/homenavi-marketplace/api/internal/db"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrListenPathInUse = errors.New("listen_path already in use")
var ErrNameInUse = errors.New("name already in use")

func ListIntegrations(ctx context.Context, db *gorm.DB, latestOnly bool, sortBy string, featuredOnly bool) ([]models.Integration, error) {
	query := db.WithContext(ctx).Model(&dbmodels.Integration{})
	if latestOnly {
		query = query.Where("latest = ?", true)
	}
	if featuredOnly {
		query = query.Where("featured = ?", true)
	}

	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "downloads":
		query = query.Order("downloads DESC, name ASC")
	case "trending":
		query = query.Order("trending_score DESC, name ASC")
	case "version":
		query = query.Order("version DESC")
	default:
		query = query.Order("name ASC")
	}

	rows := []dbmodels.Integration{}
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return mapIntegrations(rows), nil
}

func GetIntegration(ctx context.Context, db *gorm.DB, id string, version string) (*models.Integration, error) {
	query := db.WithContext(ctx).Model(&dbmodels.Integration{}).Where("id = ?", id)
	if version != "" {
		query = query.Where("version = ?", version)
	} else {
		query = query.Where("latest = ?", true)
	}
	var item dbmodels.Integration
	if err := query.First(&item).Error; err != nil {
		return nil, err
	}
	result := fromDBIntegration(item)
	return &result, nil
}

func ListVersions(ctx context.Context, db *gorm.DB, id string) ([]models.Integration, error) {
	rows := []dbmodels.Integration{}
	if err := db.WithContext(ctx).
		Model(&dbmodels.Integration{}).
		Where("id = ?", id).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return mapIntegrations(rows), nil
}

func IncrementDownloads(ctx context.Context, db *gorm.DB, id string) (*models.Integration, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&dbmodels.IntegrationDownloadEvent{IntegrationID: id}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var downloads7d int64
	if err := tx.Model(&dbmodels.IntegrationDownloadEvent{}).
		Where("integration_id = ? AND created_at >= ?", id, time.Now().Add(-7*24*time.Hour)).
		Count(&downloads7d).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&dbmodels.Integration{}).
		Where("id = ? AND latest = ?", id, true).
		Updates(map[string]any{
			"downloads":      gorm.Expr("downloads + ?", 1),
			"trending_score": float64(downloads7d),
			"updated_at":     time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var item dbmodels.Integration
	if err := tx.Where("id = ? AND latest = ?", id, true).First(&item).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	result := fromDBIntegration(item)
	return &result, nil
}

func PublishIntegration(ctx context.Context, db *gorm.DB, req models.PublishRequest, verified bool) (*models.Integration, error) {
	if req.ListenPath == "" {
		return nil, errors.New("listen_path is required")
	}

	log.Printf("store publish integration id=%q version=%q listen_path=%q verified=%t", req.ID, req.Version, req.ListenPath, verified)

	if err := ensureListenPathAvailable(ctx, db, req.ListenPath, req.ID); err != nil {
		return nil, err
	}
	if err := ensureNameAvailable(ctx, db, req.Name, req.ID); err != nil {
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

	manifestData := datatypes.JSON(manifestJSON)
	imagesData := datatypes.JSON(imagesJSON)
	assetsData := datatypes.JSON(assetsJSON)

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&dbmodels.Integration{}).Where("id = ?", req.ID).Update("latest", false).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	stats := integrationStats{Downloads: 0, TrendingScore: 0, Featured: false}
	if current, err := getLatestStats(ctx, db, req.ID); err == nil && current != nil {
		stats = *current
	}

	record := dbmodels.Integration{
		ID:            req.ID,
		Version:       req.Version,
		Name:          req.Name,
		Description:   req.Description,
		ManifestURL:   req.ManifestURL,
		Manifest:      manifestData,
		Image:         req.Image,
		Images:        imagesData,
		Assets:        assetsData,
		ListenPath:    req.ListenPath,
		ComposeFile:   req.ComposeFile,
		RepoURL:       req.RepoURL,
		ReleaseTag:    req.ReleaseTag,
		Publisher:     req.Publisher,
		Verified:      verified,
		Latest:        true,
		Downloads:     stats.Downloads,
		TrendingScore: stats.TrendingScore,
		Featured:      stats.Featured,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"description",
			"manifest_url",
			"manifest",
			"image",
			"images",
			"assets",
			"listen_path",
			"compose_file",
			"repo_url",
			"release_tag",
			"publisher",
			"verified",
			"downloads",
			"trending_score",
			"featured",
			"latest",
			"updated_at",
		}),
	}).Create(&record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	item, err := GetIntegration(ctx, db, req.ID, req.Version)
	if err != nil {
		return nil, err
	}
	log.Printf("store publish persisted id=%q version=%q latest=%t verified=%t", item.ID, item.Version, item.Latest, item.Verified)
	return item, nil
}

func ensureListenPathAvailable(ctx context.Context, db *gorm.DB, listenPath, id string) error {
	var count int64
	if err := db.WithContext(ctx).
		Model(&dbmodels.Integration{}).
		Where("listen_path = ? AND latest = ? AND id <> ?", listenPath, true, id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrListenPathInUse
	}
	return nil
}

func ensureNameAvailable(ctx context.Context, db *gorm.DB, name, id string) error {
	var count int64
	if err := db.WithContext(ctx).
		Model(&dbmodels.Integration{}).
		Where("name = ? AND latest = ? AND id <> ?", name, true, id).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrNameInUse
	}
	return nil
}

type integrationStats struct {
	Downloads     int64
	TrendingScore float64
	Featured      bool
}

func getLatestStats(ctx context.Context, db *gorm.DB, id string) (*integrationStats, error) {
	var stats integrationStats
	if err := db.WithContext(ctx).
		Model(&dbmodels.Integration{}).
		Select("downloads", "trending_score", "featured").
		Where("id = ? AND latest = ?", id, true).
		Take(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func mapIntegrations(rows []dbmodels.Integration) []models.Integration {
	out := make([]models.Integration, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBIntegration(row))
	}
	return out
}

func fromDBIntegration(row dbmodels.Integration) models.Integration {
	item := models.Integration{
		ID:          row.ID,
		Name:        row.Name,
		Version:     row.Version,
		Description: row.Description,
		ManifestURL: row.ManifestURL,
		Image:       row.Image,
		ListenPath:  row.ListenPath,
		ComposeFile: row.ComposeFile,
		RepoURL:     row.RepoURL,
		ReleaseTag:  row.ReleaseTag,
		Publisher:   row.Publisher,
		Verified:    row.Verified,
		Latest:      row.Latest,
		Downloads:   row.Downloads,
		Trending:    row.TrendingScore,
		Featured:    row.Featured,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if len(row.Manifest) > 0 {
		_ = json.Unmarshal(row.Manifest, &item.Manifest)
	}
	if len(row.Images) > 0 {
		_ = json.Unmarshal(row.Images, &item.Images)
	}
	if len(row.Assets) > 0 {
		_ = json.Unmarshal(row.Assets, &item.Assets)
	}
	return item
}
