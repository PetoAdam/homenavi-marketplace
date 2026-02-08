package models

import "time"

type Integration struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	ManifestURL string            `json:"manifest_url"`
	Manifest    map[string]any    `json:"manifest,omitempty"`
	Image       string            `json:"image"`
	Images      []string          `json:"images"`
	Assets      map[string]string `json:"assets"`
	ListenPath  string            `json:"listen_path"`
	RepoURL     string            `json:"repo_url,omitempty"`
	ReleaseTag  string            `json:"release_tag,omitempty"`
	Publisher   string            `json:"publisher,omitempty"`
	Verified    bool              `json:"verified"`
	Latest      bool              `json:"latest"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type PublishRequest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	ManifestURL string            `json:"manifest_url"`
	Manifest    map[string]any    `json:"manifest"`
	Image       string            `json:"image"`
	Images      []string          `json:"images"`
	Assets      map[string]string `json:"assets"`
	ListenPath  string            `json:"listen_path"`
	RepoURL     string            `json:"repo_url"`
	ReleaseTag  string            `json:"release_tag"`
	Publisher   string            `json:"publisher"`
}
