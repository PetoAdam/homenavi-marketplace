package handlers

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PetoAdam/homenavi-marketplace/api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type OIDCClaims struct {
	jwt.RegisteredClaims
	Repository      string `json:"repository"`
	RepositoryOwner string `json:"repository_owner"`
	Ref             string `json:"ref"`
	RefType         string `json:"ref_type"`
	SHA             string `json:"sha"`
	Workflow        string `json:"workflow"`
	JobWorkflowRef  string `json:"job_workflow_ref"`
	Actor           string `json:"actor"`
	RunID           int64  `json:"run_id"`
	RunAttempt      int64  `json:"run_attempt"`
}

type OIDCVerifier interface {
	Verify(ctx context.Context, token string) (OIDCClaims, error)
	VerifyWorkflow(ctx context.Context, claims OIDCClaims) error
}

type GitHubOIDCVerifier struct {
	issuer         string
	audience       string
	verifyWorkflow string
	githubToken    string
	client         *http.Client
	jwks           jwksCache
}

type jwksCache struct {
	mu        sync.Mutex
	expires   time.Time
	keysByKID map[string]*rsa.PublicKey
}

func NewGitHubOIDCVerifier(cfg config.Config) *GitHubOIDCVerifier {
	return &GitHubOIDCVerifier{
		issuer:         cfg.OIDCIssuer,
		audience:       cfg.OIDCAudience,
		verifyWorkflow: cfg.OIDCVerifyWorkflow,
		githubToken:    cfg.GitHubAPIToken,
		client:         &http.Client{Timeout: 10 * time.Second},
	}
}

func (v *GitHubOIDCVerifier) Verify(ctx context.Context, token string) (OIDCClaims, error) {
	var claims OIDCClaims
	if strings.TrimSpace(token) == "" {
		return claims, errors.New("missing oidc token")
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithAudience(v.audience),
		jwt.WithIssuer(v.issuer),
	)

	parsed, err := parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		key, err := v.keyForKID(ctx, kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return claims, fmt.Errorf("invalid oidc token: %w", err)
	}
	if parsed == nil || !parsed.Valid {
		return claims, errors.New("invalid oidc token")
	}

	if claims.Repository == "" || claims.Ref == "" || claims.SHA == "" {
		return claims, errors.New("missing required oidc claims")
	}

	return claims, nil
}

func (v *GitHubOIDCVerifier) VerifyWorkflow(ctx context.Context, claims OIDCClaims) error {
	if v.verifyWorkflow == "" {
		return errors.New("verify workflow not configured")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/%s/runs?per_page=1&status=success&head_sha=%s", claims.Repository, v.verifyWorkflow, claims.SHA)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "homenavi-marketplace")
	if strings.TrimSpace(v.githubToken) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(v.githubToken))
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api error: %s", resp.Status)
	}

	var payload struct {
		TotalCount int `json:"total_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if payload.TotalCount < 1 {
		return errors.New("verify workflow did not pass")
	}
	return nil
}

func (v *GitHubOIDCVerifier) keyForKID(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.jwks.mu.Lock()
	defer v.jwks.mu.Unlock()

	if v.jwks.keysByKID == nil || time.Now().After(v.jwks.expires) {
		keys, err := v.fetchJWKS(ctx)
		if err != nil {
			return nil, err
		}
		v.jwks.keysByKID = keys
		v.jwks.expires = time.Now().Add(30 * time.Minute)
	}
	key := v.jwks.keysByKID[kid]
	if key == nil {
		return nil, errors.New("unknown kid")
	}
	return key, nil
}

func (v *GitHubOIDCVerifier) fetchJWKS(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://token.actions.githubusercontent.com/.well-known/jwks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks fetch failed: %s", resp.Status)
	}

	var payload struct {
		Keys []struct {
			KID string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	out := make(map[string]*rsa.PublicKey)
	for _, key := range payload.Keys {
		if key.KID == "" || key.N == "" || key.E == "" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			continue
		}
		eInt := big.NewInt(0).SetBytes(eBytes).Int64()
		if eInt <= 0 {
			continue
		}
		out[key.KID] = &rsa.PublicKey{N: big.NewInt(0).SetBytes(nBytes), E: int(eInt)}
	}

	if len(out) == 0 {
		return nil, errors.New("no jwks keys found")
	}

	return out, nil
}
