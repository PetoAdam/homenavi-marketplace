package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/models"
	"github.com/PetoAdam/homenavi-marketplace/api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IntegrationsHandler struct {
	Pool          *pgxpool.Pool
	OIDCVerifier  OIDCVerifier
	OIDCTagPrefix string
}

func (h IntegrationsHandler) List(w http.ResponseWriter, r *http.Request) {
	latestOnly := true
	if strings.EqualFold(r.URL.Query().Get("latest"), "false") {
		latestOnly = false
	}
	items, err := store.ListIntegrations(r.Context(), h.Pool, latestOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list integrations")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"integrations": items})
}

func (h IntegrationsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	version := r.URL.Query().Get("version")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	item, err := store.GetIntegration(r.Context(), h.Pool, id, version)
	if err != nil {
		writeError(w, http.StatusNotFound, "integration not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h IntegrationsHandler) Versions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}
	items, err := store.ListVersions(r.Context(), h.Pool, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list versions")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"versions": items})
}

func (h IntegrationsHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req models.PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := validatePublishRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := store.PublishIntegration(r.Context(), h.Pool, req, true)
	if err != nil {
		if err == store.ErrListenPathInUse {
			writeError(w, http.StatusConflict, "listen_path already used")
			return
		}
		if err == store.ErrNameInUse {
			writeError(w, http.StatusConflict, "name already used")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to publish integration")
		return
	}
	log.Printf("publish stored integration id=%q version=%q latest=%t verified=%t", item.ID, item.Version, item.Latest, item.Verified)
	writeJSON(w, http.StatusOK, item)
}

func (h IntegrationsHandler) PublishOIDC(w http.ResponseWriter, r *http.Request) {
	if h.OIDCVerifier == nil {
		writeError(w, http.StatusServiceUnavailable, "oidc verifier not configured")
		return
	}

	token, err := bearerToken(r)
	if err != nil {
		log.Printf("publish-oidc unauthorized: %v", err)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := h.OIDCVerifier.Verify(r.Context(), token)
	if err != nil {
		log.Printf("publish-oidc invalid token: %v", err)
		writeError(w, http.StatusUnauthorized, "invalid oidc token")
		return
	}

	if err := h.OIDCVerifier.VerifyWorkflow(r.Context(), claims); err != nil {
		log.Printf("publish-oidc verify workflow failed: %v", err)
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	tag, err := tagFromClaims(claims, h.OIDCTagPrefix)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req models.PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("publish-oidc invalid json: %v", err)
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	log.Printf(
		"publish-oidc request fields id=%q version=%q release_tag=%q listen_path=%q repo_url=%q manifest_url=%q image=%q",
		req.ID,
		req.Version,
		req.ReleaseTag,
		req.ListenPath,
		req.RepoURL,
		req.ManifestURL,
		req.Image,
	)
	if err := validatePublishRequest(req); err != nil {
		log.Printf("publish-oidc request validation failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateOIDCRequest(req, claims, tag); err != nil {
		log.Printf("publish-oidc oidc validation failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := store.PublishIntegration(r.Context(), h.Pool, req, true)
	if err != nil {
		if err == store.ErrListenPathInUse {
			writeError(w, http.StatusConflict, "listen_path already used")
			return
		}
		if err == store.ErrNameInUse {
			writeError(w, http.StatusConflict, "name already used")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to publish integration")
		return
	}
	log.Printf("publish-oidc stored integration id=%q version=%q latest=%t verified=%t", item.ID, item.Version, item.Latest, item.Verified)
	writeJSON(w, http.StatusOK, item)
}

func validatePublishRequest(req models.PublishRequest) error {
	req.ID = strings.TrimSpace(req.ID)
	req.Name = strings.TrimSpace(req.Name)
	req.Version = strings.TrimSpace(req.Version)
	req.ListenPath = strings.TrimSpace(req.ListenPath)
	req.ManifestURL = strings.TrimSpace(req.ManifestURL)
	req.Image = strings.TrimSpace(req.Image)
	missing := make([]string, 0, 6)
	if req.ID == "" {
		missing = append(missing, "id")
	}
	if req.Name == "" {
		missing = append(missing, "name")
	}
	if req.Version == "" {
		missing = append(missing, "version")
	}
	if req.ListenPath == "" {
		missing = append(missing, "listen_path")
	}
	if req.ManifestURL == "" {
		missing = append(missing, "manifest_url")
	}
	if req.Image == "" {
		missing = append(missing, "image")
	}
	if len(missing) > 0 {
		return errField("missing required fields: " + strings.Join(missing, ", "))
	}
	if len(req.Images) > 5 {
		return errField("images must be <= 5 items")
	}
	return nil
}

func bearerToken(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	if value == "" {
		return "", errField("missing authorization header")
	}
	if !strings.HasPrefix(strings.ToLower(value), "bearer ") {
		return "", errField("invalid authorization header")
	}
	token := strings.TrimSpace(value[7:])
	if token == "" {
		return "", errField("missing bearer token")
	}
	return token, nil
}

func tagFromClaims(claims OIDCClaims, prefix string) (string, error) {
	if strings.ToLower(claims.RefType) != "tag" {
		return "", errField("oidc ref_type must be tag")
	}
	if !strings.HasPrefix(claims.Ref, "refs/tags/") {
		return "", errField("oidc ref must be a tag")
	}
	tag := strings.TrimPrefix(claims.Ref, "refs/tags/")
	if strings.TrimSpace(tag) == "" {
		return "", errField("oidc tag missing")
	}
	if strings.TrimSpace(prefix) != "" && !strings.HasPrefix(tag, prefix) {
		return "", errField("tag must start with configured prefix")
	}
	return tag, nil
}

func validateOIDCRequest(req models.PublishRequest, claims OIDCClaims, tag string) error {
	repo := strings.TrimSpace(claims.Repository)
	if repo == "" {
		log.Printf("publish-oidc missing repository claim")
		return errField("oidc repository claim missing")
	}
	log.Printf(
		"publish-oidc oidc claims repo=%q ref=%q ref_type=%q sha=%q workflow=%q job_workflow_ref=%q actor=%q",
		claims.Repository,
		claims.Ref,
		claims.RefType,
		claims.SHA,
		claims.Workflow,
		claims.JobWorkflowRef,
		claims.Actor,
	)

	if req.Version != tag {
		log.Printf("publish-oidc version mismatch: version=%q tag=%q", req.Version, tag)
		return errField("version must match the tag")
	}
	if req.ReleaseTag != tag {
		log.Printf("publish-oidc release_tag mismatch: release_tag=%q tag=%q", req.ReleaseTag, tag)
		return errField("release_tag must match the tag")
	}

	repoURL := normalizeRepoURL(req.RepoURL)
	expectedRepoURL := normalizeRepoURL("https://github.com/" + repo)
	if repoURL == "" || repoURL != expectedRepoURL {
		log.Printf("publish-oidc repo_url mismatch: repo_url=%q expected=%q", repoURL, expectedRepoURL)
		return errField("repo_url must match the GitHub repository")
	}

	rawBase := "https://raw.githubusercontent.com/" + repo + "/" + tag + "/"
	if !strings.HasPrefix(req.ManifestURL, rawBase) {
		log.Printf("publish-oidc manifest_url mismatch: manifest_url=%q expected_prefix=%q", req.ManifestURL, rawBase)
		return errField("manifest_url must point to the tag in the GitHub repo")
	}

	return nil
}

func normalizeRepoURL(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimSuffix(trimmed, ".git")
	trimmed = strings.TrimSuffix(trimmed, "/")
	return strings.ToLower(trimmed)
}

type errField string

func (e errField) Error() string {
	return string(e)
}
