package db

import (
	"time"

	"gorm.io/datatypes"
)

type Integration struct {
	ID            string `gorm:"primaryKey"`
	Version       string `gorm:"primaryKey"`
	Name          string
	Description   string
	ManifestURL   string
	Manifest      datatypes.JSON
	Image         string
	Images        datatypes.JSON
	Assets        datatypes.JSON
	ListenPath    string `gorm:"index"`
	ComposeFile   string
	RepoURL       string
	ReleaseTag    string
	Publisher     string
	Verified      bool
	Latest        bool `gorm:"index"`
	Downloads     int64
	TrendingScore float64
	Featured      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Integration) TableName() string {
	return "integrations"
}

type IntegrationDownloadEvent struct {
	ID            uint   `gorm:"primaryKey"`
	IntegrationID string `gorm:"index"`
	CreatedAt     time.Time
}

func (IntegrationDownloadEvent) TableName() string {
	return "integration_download_events"
}
