package models

import "time"

type Integration struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	ManifestURL string              `json:"manifest_url"`
	Manifest    map[string]any      `json:"manifest,omitempty"`
	Image       string              `json:"image"`
	Images      []string            `json:"images"`
	Assets      map[string]string   `json:"assets"`
	ListenPath  string              `json:"listen_path"`
	ComposeFile string              `json:"compose_file"`
	Deployment  DeploymentArtifacts `json:"deployment_artifacts"`
	RepoURL     string              `json:"repo_url,omitempty"`
	ReleaseTag  string              `json:"release_tag,omitempty"`
	Publisher   string              `json:"publisher,omitempty"`
	Verified    bool                `json:"verified"`
	Latest      bool                `json:"latest"`
	Downloads   int64               `json:"downloads"`
	Trending    float64             `json:"trending_score"`
	Featured    bool                `json:"featured"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type PublishRequest struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	ManifestURL string              `json:"manifest_url"`
	Manifest    map[string]any      `json:"manifest"`
	Image       string              `json:"image"`
	Images      []string            `json:"images"`
	Assets      map[string]string   `json:"assets"`
	ListenPath  string              `json:"listen_path"`
	ComposeFile string              `json:"compose_file"`
	Deployment  DeploymentArtifacts `json:"deployment_artifacts"`
	RepoURL     string              `json:"repo_url"`
	ReleaseTag  string              `json:"release_tag"`
	Publisher   string              `json:"publisher"`
}

type DeploymentArtifacts struct {
	Compose struct {
		File string `json:"file,omitempty"`
	} `json:"compose,omitempty"`
	Helm struct {
		ChartRef string `json:"chart_ref,omitempty"`
		Version  string `json:"version,omitempty"`
	} `json:"helm,omitempty"`
	K8sGenerated struct {
		Kind     string `json:"kind,omitempty"`
		ChartRef string `json:"chart_ref,omitempty"`
		Version  string `json:"version,omitempty"`
	} `json:"k8s_generated,omitempty"`
}
